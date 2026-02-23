package groups

import (
	"fmt"
	"slices"

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
type node struct {
	id     ID
	parent ID
	p      rl.Vector2
	level  int
	xpos   int
}

type Group struct {
	Nodes     map[ID]Node
	nodes     map[ID]node
	root      ID
	drawOrder []ID
	drawdebuc bool
}

func NewGroup(nodes map[ID]Node) *Group {
	// add missing nodes
	for _, v := range nodes {
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
	calcXPos(tree)
	calcPosiitons(tree)

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

func calcXPos(tree map[ID]node) {
	//maxW := 6 //getMaxW(tree)
	level := 0
	//prev := getLevel(tree, level)
	for {
		level++
		cur := getLevel(tree, level)
		if len(cur) == 0 {
			break
		}
		//poss := make([]int, len(cur))
		for i, id := range cur {
			n := tree[id]
			n.xpos = i
			tree[id] = n
		}
	}
}

func getLevel(tree map[ID]node, level int) []ID {
	res := []ID{}
	for k, n := range tree {
		if n.level == level {
			res = append(res, k)
		}
	}

	slices.Sort(res)

	return res
}

func calcPosiitons(tree map[ID]node) (int, int) {
	level := 0
	levels := [][]ID{}
	maxW := 0
	for {
		ns := []ID{}
		for _, n := range tree {
			if n.level == level {
				ns = append(ns, n.id)
			}
		}
		if len(ns) == 0 {
			break
		}
		maxW = max(maxW, len(ns))
		slices.SortFunc(ns, func(a, b ID) int {
			aPar := tree[tree[a].parent]
			bPar := tree[tree[b].parent]
			d := aPar.xpos - bPar.xpos
			if d != 0 {
				return d
			}
			return tree[a].xpos - tree[b].xpos
		})
		levels = append(levels, ns)

		level++
	}

	ys := 1 / float32(len(levels)-1)
	//xs := 1 / float32(maxW)
	for l := range levels {
		xs := 1 / float32(len(levels[l]))
		for c, id := range levels[l] {
			n := tree[id]

			n.p = rl.NewVector2(xs/2+xs*float32(c), ys/2+ys*float32(l-1))
			tree[n.id] = n
		}
	}

	return maxW, len(levels)
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
	g.drawdebuc = true
	return g
}

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

		s := toDraw.GetSize(rl.NewRectangle(c.X, c.Y, 100, 50))
		r := rl.NewRectangle(c.X-s.X/2, c.Y-s.Y/2, s.X, s.Y)

		rl.DrawRectangleRec(r, rl.ColorAlpha(theme.Charcoal, .2))
		rl.DrawRectangleLinesEx(r, 2, theme.Charcoal)
		g.Nodes[n.id].Draw.Draw(r)

		if g.drawdebuc {
			//rl.DrawCircleV(c, 20, theme.Salmon)
			debugTextpos := rl.Vector2Add(c, rl.NewVector2(s.X/2, -s.Y/2))
			theme.DrawText(
				fmt.Sprintf("%d: %s\n(%.2f, %.2f)\n%d", n.level, n.id, n.p.X, n.p.Y, n.xpos),
				debugTextpos, theme.Charcoal)
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
