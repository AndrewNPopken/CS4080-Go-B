package camera3d

import (
	"math"
	"image"
	"image/color"
	"../space3d"
	"../objects3d"
)

var Infinity = math.MaxFloat64

func Clamp (low, high, val * float64) float64 {
	return math.Max(*low, math.Min(*high, *val))
}

func DegreeToRadian( val float64 ) float64 {
	return val * math.Pi / 180
}

func RadianToDegree( val float64 ) float64 {
	return val * 180 / math.Pi
}

type Camera struct {
	/*
	*	In a matrix, the camera is represented as
	*
	*	X1	X2	X3	0
	*	Y1	Y2	Y3	0
	*	Z1	Z2	Z3	0
	*	T1	T2	T3	1
	*
	*	where X, Y, and Z are directional vectors pointing RIGHT, ABOVE, and BEHIND the camera respectively
	*	and T is a positional vector representing the camera's location
	*/
	ToWorld space3d.Matrix44f
}

func (c * Camera) GetPosition() space3d.Vec3f {
	return space3d.Vec3f{c.ToWorld.E[0][0], c.ToWorld.E[0][1], c.ToWorld.E[0][2]}
}

func (c * Camera) TurnLeft (rad float64) {
	cos, sin := math.Cos(rad), math.Sin(rad)
/* 	r = [][]float64{ 
		{cos, 0, sin}, 
		{0, 1, 0}, 
		{-sin, 0, cos} 
	} 
	x := c.ToWorld.E[0][0] * cos + c.ToWorld.E[0][2] * sin
	y := c.ToWorld.E[0][1]
	x := c.ToWorld.E[0][0] * (-sin) + c.ToWorld.E[0][2] * cos
*/
	c.ToWorld.E[0][0], c.ToWorld.E[0][1], c.ToWorld.E[0][2] = c.ToWorld.E[0][0] * cos + c.ToWorld.E[0][2] * sin, c.ToWorld.E[0][1], c.ToWorld.E[0][2] * cos - c.ToWorld.E[0][0] * sin
	c.ToWorld.E[1][0], c.ToWorld.E[1][1], c.ToWorld.E[1][2] = c.ToWorld.E[1][0] * cos + c.ToWorld.E[1][2] * sin, c.ToWorld.E[1][1], c.ToWorld.E[1][2] * cos - c.ToWorld.E[1][0] * sin
	c.ToWorld.E[2][0], c.ToWorld.E[2][1], c.ToWorld.E[2][2] = c.ToWorld.E[2][0] * cos + c.ToWorld.E[2][2] * sin, c.ToWorld.E[2][1], c.ToWorld.E[2][2] * cos - c.ToWorld.E[2][0] * sin
	
}

func (c * Camera) TurnUp (rad float64) {
	if c.ToWorld.E[2][1] > 0.95 && rad < 0 {
		//We are looking down
		return
	}
	if c.ToWorld.E[2][1] < -0.95 && rad > 0{
		//We are looking up
		return
	}
	cos, sin := math.Cos(rad), math.Sin(rad)
	u, w := c.ToWorld.E[0][0], c.ToWorld.E[0][2]
	
/* 	https://sites.google.com/site/glennmurray/Home/rotation-matrices-and-formulas/rotation-about-an-arbitrary-axis-in-3-dimensions
	x = (u)(u*x + w*z)(1 - cos) + x * cos - w * y * sin
	y = y * cos + ( w * x - u * z) * sin
	z = (w)(u*x + w*z)(1 - cos) + z * cos + u * y * sin
*/

	x, y, z := c.ToWorld.E[1][0], c.ToWorld.E[1][1], c.ToWorld.E[1][2]
	c.ToWorld.E[1][0], c.ToWorld.E[1][1], c.ToWorld.E[1][2] = u * (u*x + w*z) * (1 - cos) + x * cos - w * y * sin, y * cos + ( w * x - u * z) * sin, w * (u*x + w*z) * (1 - cos) + z * cos + u * y * sin

	x, y, z = c.ToWorld.E[2][0], c.ToWorld.E[2][1], c.ToWorld.E[2][2]
	c.ToWorld.E[2][0], c.ToWorld.E[2][1], c.ToWorld.E[2][2] = u * (u*x + w*z) * (1 - cos) + x * cos - w * y * sin, y * cos + ( w * x - u * z) * sin, w * (u*x + w*z) * (1 - cos) + z * cos + u * y * sin
}

type Options struct {
    Width, Height, Depth int
    FieldOfView float64
}

type Light interface {
    Light()
}

var reflectRatio float64 = 0.4
//This function takes the ray direction and turns it into a color as a placeholder. 
//Ray direction coordinates are in the range [-1,1].
//To normalized them, just add 1 and divide the result by 2.
func CastRay (origin, direction *space3d.Vec3f, objects []objects3d.Object, lights []Light, options *Options, depth int) color.RGBA {
	hitColor := color.RGBA{uint8(direction.X * 127.5 + 127.5), uint8(direction.Y * 127.5 + 127.5), uint8(direction.Z * 127.5 + 127.5), 255}
	if isHit, hitDist, hitObject := trace(origin, direction, objects); isHit {
		hitPoint := origin.Plus(direction.MultiplyBy(hitDist))
		normal := hitObject.GetSurfaceNormal(hitPoint)
		//placeholder while lights aren't implemented:
//		reflectDirection := direction.Plus(normal.MultiplyBy(direction.DotProduct(normal) * 2))
//		reflectColor := color.RGBA{uint8(reflectDirection.X * 127.5 + 127.5), uint8(reflectDirection.Y * 127.5 + 127.5), uint8(reflectDirection.Z * 127.5 + 127.5), 255}
		reflectDirection := normal
		reflectColor := color.RGBA{uint8(reflectDirection.X * 127.5 + 127.5), uint8(reflectDirection.Y * 127.5 + 127.5), uint8(reflectDirection.Z * 127.5 + 127.5), 255}
		hitColor = hitObject.GetColor()
		hitColor.R, hitColor.G, hitColor.B = uint8(float64(reflectColor.R) * reflectRatio + float64(hitColor.R) * (1 - reflectRatio)), uint8(float64(reflectColor.G) * reflectRatio + float64(hitColor.G) * (1 - reflectRatio)), uint8(float64(reflectColor.B) * reflectRatio + float64(hitColor.B) * (1 - reflectRatio))
		/*
		////This may be important later if we add textures
		// Use the normal and texture coordinates to shade the hit point.
		// The normal is used to compute a simple facing ratio and the texture coordinate
		// to compute a basic checker board pattern
		float scale = 4
		float pattern = (fmodf(tex.x * scale, 1) > 0.5) ^ (fmodf(tex.y * scale, 1) > 0.5)
		hitColor = std::max(0.f, Nhit.dotProduct(-dir)) * mix(hitObject->color, hitObject->color * 0.8, pattern)
		*/
	}
	return hitColor
}

//This function iterates over all objects to find the first visible intersection
//returns bool, float64, objects3d.Object
//bool is whether an object was hit
//float64 is distance to the intersection point
//objects3d.Object is the object hit
func trace (origin, direction *space3d.Vec3f, objects []objects3d.Object) (bool, float64, objects3d.Object) {
	dist := Infinity
	var hitObject objects3d.Object
	for i := 0; i < len(objects); i++ {
		if hit, tempDist := objects[i].Intersect(origin, direction); hit && tempDist < dist {
			hitObject = objects[i]
			dist = tempDist
		}
	}
	return hitObject != nil, dist, hitObject
}

//The main render function. This where we iterate over all pixels in the image, generate
//primary rays and cast these rays into the scene. The content of the framebuffer is returned as an RGBA image
func Render (camera * Camera, objects []objects3d.Object, lights []Light, options *Options) *image.RGBA {
	//CameraToWorld := space3d.NewIdentityMatrix()
	framebuffer := image.NewRGBA(image.Rect(0, 0, options.Width, options.Height))
	scale := math.Tan(DegreeToRadian(options.FieldOfView * 0.5))
	aspectRatio := float64(options.Width) / float64(options.Height)
	//Transform the camera origin by transforming the point with coordinates (0,0,0) 
	//to world-space using the camera.ToWorld matrix.
	origin := camera.GetPosition()
	camera.ToWorld.MultiplyVectorMatrix(space3d.Vec3f{}, &origin)
	//Linear for now, d==down, r==right
	for d := 0; d < options.Height; d++ {
		for r := 0; r < options.Width; r++ {
			x := (float64(2 * (r - options.Width / 2) + 1) / float64(options.Width - 1)) * scale
			y := (1 - float64(2 * d + 1) / float64(options.Height)) * scale / aspectRatio
			var direction space3d.Vec3f
			camera.ToWorld.MultiplyDirectionalMatrix(space3d.Vec3f{x, y, -1}, &direction)
			direction.Normalize()
			framebuffer.SetRGBA(r, d, CastRay(&origin, &direction, objects, lights, options, options.Depth))
		}
	}
	return framebuffer
}
