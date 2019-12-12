package avl

import "testing"

func benchmarkTreeGeneratorN(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		tr := NewTreeGenerator()
		for j := 1; j <= n; j++ {
			if _, err := tr.Add(newExampleMutableNode(j)); err != nil {
				panic(err)
			}
		}
	}
}

func BenchmarkTreeGenerator10(b *testing.B)    { benchmarkTreeGeneratorN(10, b) }
func BenchmarkTreeGenerator100(b *testing.B)   { benchmarkTreeGeneratorN(100, b) }
func BenchmarkTreeGenerator200(b *testing.B)   { benchmarkTreeGeneratorN(200, b) }
func BenchmarkTreeGenerator300(b *testing.B)   { benchmarkTreeGeneratorN(300, b) }
func BenchmarkTreeGenerator400(b *testing.B)   { benchmarkTreeGeneratorN(400, b) }
func BenchmarkTreeGenerator500(b *testing.B)   { benchmarkTreeGeneratorN(500, b) }
func BenchmarkTreeGenerator600(b *testing.B)   { benchmarkTreeGeneratorN(600, b) }
func BenchmarkTreeGenerator700(b *testing.B)   { benchmarkTreeGeneratorN(700, b) }
func BenchmarkTreeGenerator800(b *testing.B)   { benchmarkTreeGeneratorN(800, b) }
func BenchmarkTreeGenerator900(b *testing.B)   { benchmarkTreeGeneratorN(900, b) }
func BenchmarkTreeGenerator1000(b *testing.B)  { benchmarkTreeGeneratorN(1000, b) }
func BenchmarkTreeGenerator1100(b *testing.B)  { benchmarkTreeGeneratorN(1100, b) }
func BenchmarkTreeGenerator1200(b *testing.B)  { benchmarkTreeGeneratorN(1200, b) }
func BenchmarkTreeGenerator1300(b *testing.B)  { benchmarkTreeGeneratorN(1300, b) }
func BenchmarkTreeGenerator1400(b *testing.B)  { benchmarkTreeGeneratorN(1400, b) }
func BenchmarkTreeGenerator1500(b *testing.B)  { benchmarkTreeGeneratorN(1500, b) }
func BenchmarkTreeGenerator1600(b *testing.B)  { benchmarkTreeGeneratorN(1600, b) }
func BenchmarkTreeGenerator1700(b *testing.B)  { benchmarkTreeGeneratorN(1700, b) }
func BenchmarkTreeGenerator1800(b *testing.B)  { benchmarkTreeGeneratorN(1800, b) }
func BenchmarkTreeGenerator1900(b *testing.B)  { benchmarkTreeGeneratorN(1900, b) }
func BenchmarkTreeGenerator10000(b *testing.B) { benchmarkTreeGeneratorN(10000, b) }
