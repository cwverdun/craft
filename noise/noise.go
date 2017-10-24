package noise

import "math"

const (
	f2 = 0.3660254037844386
	g2 = 0.21132486540518713
	f3 = 1.0 / 3.0
	g3 = 1.0 / 6.0
)

func assign(a [3]int, v0, v1, v2 int) {
	a[0] = v0
	a[1] = v1
	a[2] = v2
}

func dot3(v1, v2 [3]float64) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

var grad3 = [16][3]float64{
	{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
	{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
	{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1},
	{1, 0, -1}, {-1, 0, -1}, {0, -1, 1}, {0, 1, 1},
}

var perm = []int{
	151, 160, 137, 91, 90, 15, 131, 13,
	201, 95, 96, 53, 194, 233, 7, 225,
	140, 36, 103, 30, 69, 142, 8, 99,
	37, 240, 21, 10, 23, 190, 6, 148,
	247, 120, 234, 75, 0, 26, 197, 62,
	94, 252, 219, 203, 117, 35, 11, 32,
	57, 177, 33, 88, 237, 149, 56, 87,
	174, 20, 125, 136, 171, 168, 68, 175,
	74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122,
	60, 211, 133, 230, 220, 105, 92, 41,
	55, 46, 245, 40, 244, 102, 143, 54,
	65, 25, 63, 161, 1, 216, 80, 73,
	209, 76, 132, 187, 208, 89, 18, 169,
	200, 196, 135, 130, 116, 188, 159, 86,
	164, 100, 109, 198, 173, 186, 3, 64,
	52, 217, 226, 250, 124, 123, 5, 202,
	38, 147, 118, 126, 255, 82, 85, 212,
	207, 206, 59, 227, 47, 16, 58, 17,
	182, 189, 28, 42, 223, 183, 170, 213,
	119, 248, 152, 2, 44, 154, 163, 70,
	221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110,
	79, 113, 224, 232, 178, 185, 112, 104,
	218, 246, 97, 228, 251, 34, 242, 193,
	238, 210, 144, 12, 191, 179, 162, 241,
	81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157,
	184, 84, 204, 176, 115, 121, 50, 45,
	127, 4, 150, 254, 138, 236, 205, 93,
	222, 114, 67, 29, 24, 72, 243, 141,
	128, 195, 78, 66, 215, 61, 156, 180,
	151, 160, 137, 91, 90, 15, 131, 13,
	201, 95, 96, 53, 194, 233, 7, 225,
	140, 36, 103, 30, 69, 142, 8, 99,
	37, 240, 21, 10, 23, 190, 6, 148,
	247, 120, 234, 75, 0, 26, 197, 62,
	94, 252, 219, 203, 117, 35, 11, 32,
	57, 177, 33, 88, 237, 149, 56, 87,
	174, 20, 125, 136, 171, 168, 68, 175,
	74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122,
	60, 211, 133, 230, 220, 105, 92, 41,
	55, 46, 245, 40, 244, 102, 143, 54,
	65, 25, 63, 161, 1, 216, 80, 73,
	209, 76, 132, 187, 208, 89, 18, 169,
	200, 196, 135, 130, 116, 188, 159, 86,
	164, 100, 109, 198, 173, 186, 3, 64,
	52, 217, 226, 250, 124, 123, 5, 202,
	38, 147, 118, 126, 255, 82, 85, 212,
	207, 206, 59, 227, 47, 16, 58, 17,
	182, 189, 28, 42, 223, 183, 170, 213,
	119, 248, 152, 2, 44, 154, 163, 70,
	221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110,
	79, 113, 224, 232, 178, 185, 112, 104,
	218, 246, 97, 228, 251, 34, 242, 193,
	238, 210, 144, 12, 191, 179, 162, 241,
	81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157,
	184, 84, 204, 176, 115, 121, 50, 45,
	127, 4, 150, 254, 138, 236, 205, 93,
	222, 114, 67, 29, 24, 72, 243, 141,
	128, 195, 78, 66, 215, 61, 156, 180,
}

//Simplex2 ...
func Simplex2(x, y float64, octaves int, persistence, lacunarity float64) float64 {
	freq, amp, max := 1.0, 1.0, 1.0
	total := noise2(x, y)
	for i := 1; i < octaves; i++ {
		freq *= lacunarity
		amp *= persistence
		max += amp
		total += noise2(x*freq, y*freq) * amp
	}
	return (1 + total/max) / 2
}

func noise2(x, y float64) float64 {
	var i1, j1, I, J int
	s := (x + y) * f2
	i := math.Floor(x + s)
	j := math.Floor(y + s)
	t := (i + j) * g2

	var xx, yy, f, noise [3]float64
	var g [3]int

	xx[0] = x - (i - t)
	yy[0] = y - (j - t)

	if xx[0] > yy[0] {
		i1 = 1
	}
	if xx[0] <= yy[0] {
		j1 = 1
	}

	xx[2] = xx[0] + g2*2.0 - 1.0
	yy[2] = yy[0] + g2*2.0 - 1.0
	xx[1] = xx[0] - float64(i1) + g2
	yy[1] = yy[0] - float64(j1) + g2

	I = int(i) & 255
	J = int(j) & 255
	g[0] = perm[I+perm[J]] % 12
	g[1] = perm[I+i1+perm[J+j1]] % 12
	g[2] = perm[I+1+perm[J+1]] % 12

	for c := 0; c <= 2; c++ {
		f[c] = 0.5 - xx[c]*xx[c] - yy[c]*yy[c]
	}

	for c := 0; c <= 2; c++ {
		if f[c] > 0 {
			noise[c] = f[c] * f[c] * f[c] * f[c] * (grad3[g[c]][0]*xx[c] + grad3[g[c]][1]*yy[c])
		}
	}

	return (noise[0] + noise[1] + noise[2]) * 70.0
}

//Simplex3 ...
func Simplex3(x, y, z float64, octaves int, persistence, lacunarity float64) float64 {
	freq, amp, max := 1.0, 1.0, 1.0
	total := noise2(x, y)
	//++i was here TODO
	for i := 1; i < octaves; i++ {
		freq *= lacunarity
		amp *= persistence
		max += amp
		total += noise3(x*freq, y*freq, z*freq) * amp
	}
	return (1 + total/max) / 2
}

func noise3(x, y, z float64) float64 {
	var o1, o2 [3]int
	var g [4]int
	var f, noise [4]float64
	s := (x + y + z) * f3
	i := math.Floor(x + s)
	j := math.Floor(y + s)
	k := math.Floor(z + s)
	t := (i + j + k) * g3

	var pos [4][3]float64

	pos[0][0] = x - (i - t)
	pos[0][1] = y - (j - t)
	pos[0][2] = z - (k - t)

	if pos[0][0] >= pos[0][1] {
		if pos[0][1] >= pos[0][2] {
			assign(o1, 1, 0, 0)
			assign(o2, 1, 1, 0)
		} else if pos[0][0] >= pos[0][2] {
			assign(o1, 1, 0, 0)
			assign(o2, 1, 0, 1)
		} else {
			assign(o1, 0, 0, 1)
			assign(o2, 1, 0, 1)
		}
	} else {
		if pos[0][1] < pos[0][2] {
			assign(o1, 0, 0, 1)
			assign(o2, 0, 1, 1)
		} else if pos[0][0] < pos[0][2] {
			assign(o1, 0, 1, 0)
			assign(o2, 0, 1, 1)
		} else {
			assign(o1, 0, 1, 0)
			assign(o2, 1, 1, 0)
		}
	}

	for c := 0; c <= 2; c++ {
		pos[3][c] = pos[0][c] - 1.0 + 3.0*g3
		pos[2][c] = pos[0][c] - float64(o2[c]) + 2.0*g3
		pos[1][c] = pos[0][c] - float64(o1[c]) + g3
	}

	I := int(i) & 255
	J := int(j) & 255
	K := int(k) & 255
	g[0] = perm[I+perm[J+perm[K]]] % 12
	g[1] = perm[I+o1[0]+perm[J+o1[1]+perm[o1[2]+K]]] % 12
	g[2] = perm[I+o2[0]+perm[J+o2[1]+perm[o2[2]+K]]] % 12
	g[3] = perm[I+1+perm[J+1+perm[K+1]]] % 12

	for c := 0; c <= 3; c++ {
		f[c] = 0.6 - pos[c][0]*pos[c][0] - pos[c][1]*pos[c][1] - pos[c][2]*pos[c][2]
	}

	for c := 0; c <= 3; c++ {
		if f[c] > 0 {
			noise[c] = f[c] * f[c] * f[c] * f[c] * dot3(pos[c], grad3[g[c]])
		}
	}

	return (noise[0] + noise[1] + noise[2] + noise[3]) * 32.0
}