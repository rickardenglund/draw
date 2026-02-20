#version 330

in vec2 fragTexCoord;
in vec4 fragColor;
in vec2 fragPos;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;

uniform float seconds;
#define MAX_PS  400
uniform vec2 ps[MAX_PS];

// Output fragment color
out vec4 finalColor;


uniform vec2 center = vec2(.5, .5);

float DistToLine(vec2 a, vec2 b, vec2 p)
{
    vec2 ba = b - a;
    vec2 pa = p-a;
    float h = clamp(dot(pa, ba)/dot(ba, ba), 0, 1);
    return length(pa-h*ba);
}

void main()
{
    float m = 1;

//    for (int i=0; i<MAX_PS-1; i++)
//    {
//        m = min(m, DistToLine(ps[i], ps[i+1], fragPos));
//    }


//    float frame_width = .02;
//    float frameOpacity = smoothstep(frame_width, 0, fragTexCoord.x) + smoothstep(1-frame_width, 1, fragTexCoord.x) +
//         smoothstep(frame_width, 0, fragTexCoord.y) + smoothstep(1-frame_width, 1, fragTexCoord.y);

    vec4 graph = vec4(fragColor.rgb, 1);

//    float graph_width = .0005;
//    float graph_step_w = .001;
//    float grap_opacity = smoothstep(graph_width+ graph_step_w, graph_width, m);

//    vec3 frameClr = vec3(.2, .2, .2);
//    finalColor = vec4(1,0,0,1);//graph*grap_opacity + vec4(frameClr, 1)*frameOpacity;
    finalColor = vec4(fragPos.xy, 0, 1);
}

