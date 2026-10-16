package groups

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
	"github.com/rickardenglund/draw/views/text"
)

type ID string
type Node struct {
	Parent ID
	Draw   draw.Drawable
}
// Frameless is implemented by node drawables that draw their own shape, e.g.
// an icon with a label, instead of getting the default box.
type Frameless interface {
	Frameless()
}

type node struct {
	id     ID
	parent ID
	p      rl.Vector2
	level  int
	xpos   int
	// slot is the share of the group's width and height the node may use.
	slot rl.Vector2
}

type Group struct {
	Nodes     map[ID]Node
	nodes     map[ID]node
	root      ID
	drawOrder []ID
	drawDebug bool
}

func NewGroup(nodes map[ID]Node) *Group {
	// add missing nodes
	for k, v := range nodes {
		if k == "" {
			delete(nodes, k)
			newID := genID()
			nodes[newID] = v
		}
		_, exist := nodes[v.Parent]
		if !exist && v.Parent != "" {
			nodes[v.Parent] = Node{
				Parent: "",
				Draw:   nil,
			}
		}
	}

	rootID := ID("root")
	nodes[rootID] = Node{
		Parent: "",
		Draw:   nil,
	}

	tree := map[ID]node{}
	for k, v := range nodes {
		if v.Draw == nil {
			v.Draw = text.NewTextf("%s", k)
			nodes[k] = v
		}
		parent := v.Parent
		if parent == "" && k != rootID {
			parent = rootID
		}
		tree[k] = node{
			id:     k,
			parent: parent,
			level:  0,
		}
	}
	calcLevels(tree, rootID)
	calcPositions(tree, rootID)

	drawlist := make([]ID, 0, len(tree))
	for _, n := range tree {
		if n.id == rootID {
			continue
		}
		drawlist = append(drawlist, n.id)
	}
	slices.Sort(drawlist)

	return &Group{
		Nodes:     nodes,
		drawOrder: drawlist,
		nodes:     tree,
		root:      rootID,
	}
}

var nextID atomic.Int32

func genID() ID {
	return ID(fmt.Sprintf("gen-%d", nextID.Add(1)))
}

// calcPositions lays the tree out top-down: every subtree gets a horizontal
// share proportional to its number of leaves, children are ordered by ID and
// placed under their parent, and parents are centred over their children.
func calcPositions(tree map[ID]node, rootID ID) {
	children := map[ID][]ID{}
	depth := 0
	for id, n := range tree {
		depth = max(depth, n.level)
		if id != rootID {
			children[n.parent] = append(children[n.parent], id)
		}
	}
	for _, cs := range children {
		slices.SortFunc(cs, naturalCompare)
	}

	leaves := map[ID]int{}
	var countLeaves func(ID) int
	countLeaves = func(id ID) int {
		cs := children[id]
		if len(cs) == 0 {
			leaves[id] = 1
			return 1
		}
		total := 0
		for _, c := range cs {
			total += countLeaves(c)
		}
		leaves[id] = total
		return total
	}
	total := float32(countLeaves(rootID))

	ys := float32(1)
	if depth > 0 {
		ys = 1 / float32(depth)
	}
	var place func(id ID, offset int)
	place = func(id ID, offset int) {
		n := tree[id]
		n.xpos = offset
		xs := float32(leaves[id]) / total
		n.p = rl.NewVector2(float32(offset)/total+xs/2, ys/2+ys*float32(n.level-1))
		n.slot = rl.NewVector2(xs, ys)
		tree[id] = n

		for _, c := range children[id] {
			place(c, offset)
			offset += leaves[c]
		}
	}
	place(rootID, 0)
}

func calcLevels(tree map[ID]node, rootID ID) {
	toVisit := []ID{rootID}
	visited := map[ID]bool{}

	for len(toVisit) > 0 {
		cur := tree[toVisit[0]]
		visited[cur.id] = true
		toVisit = toVisit[1:]

		cs := getChildres(tree, cur.id)
		for _, c := range cs {
			tree[c] = node{
				id:     c,
				parent: cur.id,
				level:  cur.level + 1,
			}

			toVisit = append(toVisit, c)
			if visited[c] {
				panic("loop")
			}
		}
		visited[cur.id] = true

	}
}

func getChildres(tree map[ID]node, id ID) []ID {
	cs := []ID{}
	for _, v := range tree {
		if v.parent == id {
			cs = append(cs, v.id)
		}
	}
	return cs
}

func (g *Group) WithDebug() *Group {
	g.drawDebug = true
	return g
}

const (
	nodeMargin  = 20
	minNodeSize = 10
)

func (g *Group) Draw(t rl.Rectangle) {
	for _, id := range g.drawOrder {
		n := g.nodes[id]
		if n.parent == g.root {
			continue
		}

		par, ok := g.nodes[n.parent]
		if !ok {
			continue
		}
		pPos := pToPos(t, par.p)
		nPos := pToPos(t, n.p)
		rl.DrawLineEx(nPos, pPos, 2, rl.ColorBrightness(theme.Charcoal, .3))
		rl.DrawLineEx(nPos, pPos, 1, rl.ColorBrightness(theme.Charcoal, .5))
	}

	for _, id := range g.drawOrder {
		n := g.nodes[id]
		c := pToPos(t, n.p)
		toDraw := g.Nodes[n.id].Draw

		// Nodes take their natural size, bounded by their slot so neighbours don't overlap.
		bound := rl.NewRectangle(c.X, c.Y,
			max(0, t.Width*n.slot.X-nodeMargin),
			max(0, t.Height*n.slot.Y-nodeMargin))
		s := toDraw.GetSize(bound)
		// Keep a visible box even when the content doesn't fit at all.
		s = rl.NewVector2(max(s.X, minNodeSize), max(s.Y, minNodeSize))
		r := rl.NewRectangle(c.X-s.X/2, c.Y-s.Y/2, s.X, s.Y)

		if _, ok := toDraw.(Frameless); !ok {
			// Opaque backing so edges don't show through the node.
			rl.DrawRectangleRec(r, theme.Nude)
			rl.DrawRectangleRec(r, rl.ColorAlpha(theme.Charcoal, .2))
			rl.DrawRectangleLinesEx(r, 2, theme.Charcoal)
		}
		toDraw.Draw(r)

		if g.drawDebug {
			//rl.DrawCircleV(c, 20, theme.Salmon)
			debugTextPos := rl.Vector2Add(c, rl.NewVector2(s.X/2, -s.Y/2))
			theme.DrawText(
				fmt.Sprintf("%d: %s\n(%.2f, %.2f)\n%d", n.level, n.id, n.p.X, n.p.Y, n.xpos),
				debugTextPos, theme.Charcoal)
		}
	}

}

func pToPos(t rl.Rectangle, p rl.Vector2) rl.Vector2 {
	start := rl.NewVector2(t.X, t.Y)
	c := rl.Vector2Add(start, rl.NewVector2(t.Width*p.X, t.Height*p.Y))
	return c
}

func (g *Group) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.Width, target.Height)
}

// naturalCompare orders runs of digits by value, so v2 comes before v10.
func naturalCompare(a, b ID) int {
	for a != "" && b != "" {
		da, db := digitPrefix(a), digitPrefix(b)
		if da != "" && db != "" {
			na, _ := strconv.Atoi(string(da))
			nb, _ := strconv.Atoi(string(db))
			if c := cmp.Compare(na, nb); c != 0 {
				return c
			}
			a, b = a[len(da):], b[len(db):]
			continue
		}
		if c := cmp.Compare(a[0], b[0]); c != 0 {
			return c
		}
		a, b = a[1:], b[1:]
	}
	return cmp.Compare(len(a), len(b))
}

func digitPrefix(s ID) ID {
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	return s[:i]
}
