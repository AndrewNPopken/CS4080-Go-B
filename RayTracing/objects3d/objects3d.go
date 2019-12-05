package objects3d

import (
	"../space3d"
	"image/color"
	"math"
)

var Infinity = math.MaxFloat64

type Object interface {
	//bool is if there is an intersection, float64 is distance from origin to point hit (if hit)
	Intersect(origin, direction * space3d.Vec3f) (bool, float64)
	GetColor() color.RGBA
	GetSurfaceNormal(pointOfIntersection space3d.Vec3f) space3d.Vec3f
}

type Sphere struct {
	Position space3d.Vec3f
	Radius float64
	Color color.RGBA
}

/*
func (Object ob) BoundingSphere() Sphere {
	var bounds Sphere
	//bounds.Position = //rough center of ob //objects either have a recorded center or the means to estimate
	bounds.Position = ob.Position()
	//bounds.Radius = //at least the distance from ob's Position to the furthest point in ob's geometry //default to 1
	bounds.Radius = 1.0
	return bounds
}
*/


func (ob Sphere) Intersect (origin, direction * space3d.Vec3f) (bool, float64) {
	cameraToObject := ob.Position.Minus(*origin)
	camDirDot := cameraToObject.DotProduct(*direction)
	if camDirDot < 0 {
	//sphere is behind camera
		return false, Infinity
	}
	magnitudeDifference := cameraToObject.DotProduct(cameraToObject) - camDirDot * camDirDot
	radiusSquared := ob.Radius * ob.Radius
	if magnitudeDifference > radiusSquared {
	//ray misses the sphere
		return false, Infinity
	}
	distMidpoint := math.Sqrt(radiusSquared - magnitudeDifference)
	dist1 := camDirDot - distMidpoint
	dist2 := camDirDot + distMidpoint

	if dist1 < dist2 && dist1 > 0{
		//dist1 is distance to closest intersection on sphere
		return true, dist1
	} else if dist2 > 0 {
		//dist2 is distance to closest intersection on sphere, somehow (this is just in case we clip the sphere)
		return true, dist2
	}
	//ending up here should be impossible
	//cause we confirmed the sphere is in front of the camera
	//but to pass the above if/else if chain without a true implies that the camera is in front of the sphere
	//but just in case...
	return false, Infinity
}

func (ob Sphere) GetColor() color.RGBA {
	return ob.Color
}

func (ob Sphere) GetSurfaceNormal(pointOfIntersection space3d.Vec3f) space3d.Vec3f {
	temp := pointOfIntersection.Minus(ob.Position)
	temp.Normalize()
	return temp
}