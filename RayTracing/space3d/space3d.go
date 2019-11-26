package space3d

import "math"
//import "fmt"

type Vec3f struct {
	X, Y, Z float64
}

func (v Vec3f) Plus (o Vec3f) Vec3f {
	return Vec3f{ X: v.X + o.X, Y: v.Y + o.Y, Z: v.Z + o.Z }
}

func (v * Vec3f) AddAndSet (o Vec3f) {
	v.X, v.Y, v.Z = v.X + o.X, v.Y + o.Y, v.Z + o.Z
}

func (v Vec3f) Minus (o Vec3f) Vec3f {
	return Vec3f{ X: v.X - o.X, Y: v.Y - o.Y, Z: v.Z - o.Z }
}

func (v * Vec3f) SubtractAndSet (o Vec3f) {
	v.X, v.Y, v.Z = v.X - o.X, v.Y - o.Y, v.Z - o.Z
}

func (v Vec3f) MultiplyBy (o float64) Vec3f {
	return Vec3f{ X: v.X * o, Y: v.Y * o, Z: v.Z * o }
}

func (v * Vec3f) MultiplyAndSet (o float64) {
	v.X, v.Y, v.Z = v.X * o, v.Y * o, v.Z * o
}

func (v Vec3f) DivideBy (o float64) Vec3f {
	return Vec3f{ X: v.X / o, Y: v.Y / o, Z: v.Z / o }
}

func (v * Vec3f) DivideAndSet (o float64) {
	v.X, v.Y, v.Z = v.X / o, v.Y / o, v.Z / o
}

func (v Vec3f) CrossProduct (o Vec3f) Vec3f {
	return Vec3f{ X: v.Y*o.Z - v.Z*o.Y, Y: v.X*o.Z - v.Z*o.X, Z: v.X*o.Y - v.Y*o.X }
}

func (v Vec3f) DotProduct (o Vec3f) float64 {
	return v.X * o.X + v.Y * o.Y + v.Z * o.Z 
}

func (v Vec3f) Negative () Vec3f {
	return Vec3f{ X: - v.X, Y: - v.Y, Z: - v.Z }
}

func (v Vec3f) Norm () float64 {
	return v.X * v.X + v.Y * v.Y + v.Z * v.Z
}

func (v Vec3f) Length () float64 {
	return math.Sqrt(v.X * v.X + v.Y * v.Y + v.Z * v.Z)
}

func (v * Vec3f) Normalize () {
	norm := v.Norm()
	if norm > 0 {
		length := math.Sqrt(norm)
		v.X, v.Y, v.Z = v.X / length, v.Y / length, v.Z / length
	}
}

type Matrix44f struct {
	E[4][4] float64
}

func NewIdentityMatrix () Matrix44f {
	return Matrix44f{E:[4][4]float64{ {1,0,0,0}, {0,1,0,0}, {0,0,1,0}, {0,0,0,1} }}
}

func NewEmptyMatrix () Matrix44f {
	return Matrix44f{}
}

func NewDefinedMatrix (m [4][4]float64) Matrix44f {
	return Matrix44f{E: m }
}

func (m Matrix44f) MultiplyBy (o Matrix44f) Matrix44f {
	temp := Matrix44f{}
	for i := 0; i < 4; i++ { 
		for j := 0; j < 4; j++ { 
			temp.E[i][j] = m.E[i][0] * o.E[0][j] + m.E[i][1] * o.E[1][j] + m.E[i][2] * o.E[2][j] + m.E[i][3] * o.E[3][j]
		}
	}
	return temp
}

func (m Matrix44f) TransposedCopy () Matrix44f {
	return Matrix44f{E:[4][4]float64{ {m.E[0][0], m.E[1][0], m.E[2][0], m.E[3][0]}, {m.E[0][1], m.E[1][1], m.E[2][1], m.E[3][1]}, {m.E[0][2], m.E[1][2], m.E[2][2], m.E[3][2]}, {m.E[0][3], m.E[1][3], m.E[2][3], m.E[3][3]} }}
}

func (m * Matrix44f) TransposeSelf () {
	m.E[0][1], m.E[0][2], m.E[0][3], m.E[1][0], m.E[1][2], m.E[1][3], m.E[2][0], m.E[2][1], m.E[2][3], m.E[3][0], m.E[3][1], m.E[3][2] = m.E[1][0], m.E[2][0], m.E[3][0], m.E[0][1], m.E[2][1], m.E[3][1], m.E[0][2], m.E[1][2], m.E[3][2], m.E[0][3], m.E[1][3], m.E[2][3]
}

//void multVecMatrix(const Vec3<S> &src, Vec3<S> &dst)
func (m Matrix44f) MultiplyVectorMatrix (origin Vec3f, destination * Vec3f) {
	x := origin.X * m.E[0][0] + origin.Y * m.E[1][0] + origin.Z * m.E[2][0] + m.E[3][0]
	y := origin.X * m.E[0][1] + origin.Y * m.E[1][1] + origin.Z * m.E[2][1] + m.E[3][1]
	z := origin.X * m.E[0][2] + origin.Y * m.E[1][2] + origin.Z * m.E[2][2] + m.E[3][2]
	w := origin.X * m.E[0][3] + origin.Y * m.E[1][3] + origin.Z * m.E[2][3] + m.E[3][3]

	destination.X, destination.Y, destination.Z = x / w, y / w, z / w
}

//void multDirMatrix(const Vec3<S> &src, Vec3<S> &dst) const
func (m Matrix44f) MultiplyDirectionalMatrix (origin Vec3f, destination * Vec3f) {
	x := origin.X * m.E[0][0] + origin.Y * m.E[1][0] + origin.Z * m.E[2][0]
	y := origin.X * m.E[0][1] + origin.Y * m.E[1][1] + origin.Z * m.E[2][1]
	z := origin.X * m.E[0][2] + origin.Y * m.E[1][2] + origin.Z * m.E[2][2]

	destination.X, destination.Y, destination.Z = x, y, z
}

//Matrix44 inverse() const 
func (m Matrix44f) Inverse () Matrix44f { 
	var i, j, k int
	result := NewIdentityMatrix()
	mcopy := NewDefinedMatrix(m.E)
	//fmt.Println(result)
	//fmt.Println(mcopy)

	//Forward elimination
	for i = 0; i < 3; i++ {
		pivot := i
		pivotsize := mcopy.E[i][i]
		if pivotsize < 0 {
			pivotsize = -pivotsize
		}
		for j = i + 1; j < 4; j++ {
			temp := mcopy.E[j][i]
			if temp < 0 {
				temp = -temp
			}
			if temp > pivotsize {
				pivot, pivotsize = j, temp
			}
		}
		if pivotsize == 0 {
			// Cannot invert singular matrix
			return NewIdentityMatrix()
		}
		if pivot != i {
			for j = 0; j < 4; j++ {
				mcopy.E[i][j], mcopy.E[pivot][j], result.E[i][j], result.E[pivot][j] = mcopy.E[pivot][j], mcopy.E[i][j], result.E[pivot][j], result.E[i][j]
			}
		}
		for j = i + 1; j < 4; j++ {
			temp := mcopy.E[j][i] / mcopy.E[i][i]

			for k = 0; k < 4; k++ {
				mcopy.E[j][k] -= temp * mcopy.E[i][k]
				result.E[j][k] -= temp * result.E[i][k]
			}
		}
	}
	// Backward substitution
	for i = 3; i >= 0; i-- {
		temp := mcopy.E[i][i]
		if temp == 0 {
			// Cannot invert singular matrix
			return NewIdentityMatrix()
		}
		for j = 0; j < 4; j++ {
			mcopy.E[i][j] /= temp
			result.E[i][j] /= temp
		}
		for j = 0; j < i; j++ {
			temp = mcopy.E[j][i]
			for k = 0; k < 4; k++ {
				mcopy.E[j][k] -= temp * mcopy.E[i][k]
				result.E[j][k] -= temp * result.E[i][k]
			}
		}
	}
	//fmt.Println(result)
	return result
}

func (m * Matrix44f) Invert () Matrix44f{
	*m = m.Inverse()
	return *m
}

