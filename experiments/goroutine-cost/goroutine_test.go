package goroutinecost

import "testing"

func BenchmarkSpawnUnbuffered(b *testing.B) {
	for b.Loop() {
		SpawnUnbuffered()
	}
}

func BenchmarkSpawnBuffered(b *testing.B) {
	for b.Loop() {
		SpawnBuffered()
	}
}

func BenchmarkSpawnWaitGroup(b *testing.B) {
	for b.Loop() {
		SpawnWaitGroup()
	}
}
