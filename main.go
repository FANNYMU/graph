package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Konfigurasi optimasi ekstrem
const (
	CONCURRENT_WORKERS = 16   // Jumlah worker goroutine paralel
	BATCH_SIZE         = 1000 // Jumlah commit dalam satu batch
)

func main() {
	fmt.Print("\033[36m=== GITHUB FAKE COMMIT GENERATOR TURBO ===\033[0m\n")
	fmt.Print("\033[33mMasukkan jumlah commit yang diinginkan: \033[0m")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\033[31mError: %s\033[0m\n", err)
		return
	}

	input = strings.TrimSpace(input)
	targetCommits, err := strconv.Atoi(input)
	if err != nil {
		fmt.Print("\033[31mError: Masukkan angka yang valid!\033[0m\n")
		return
	}

	// Konfigurasi Git untuk performa maksimal
	setupGitForPerformance()

	buatFakeCommit(targetCommits)
}

func setupGitForPerformance() {
	// Konfigurasi Git untuk performa maksimal
	commands := [][]string{
		{"git", "config", "--local", "gc.auto", "0"},                    // Matikan auto garbage collection
		{"git", "config", "--local", "gc.autoDetach", "false"},          // Matikan auto detach
		{"git", "config", "--local", "core.fsyncObjectFiles", "false"},  // Matikan fsync
		{"git", "config", "--local", "core.preloadIndex", "true"},       // Aktifkan preload index
		{"git", "config", "--local", "core.commitGraph", "true"},        // Aktifkan commit graph
		{"git", "config", "--local", "feature.manyFiles", "true"},       // Optimasi untuk banyak file
		{"git", "config", "--local", "advice.detachedHead", "false"},    // Matikan warning detached head
		{"git", "config", "--local", "status.showUntrackedFiles", "no"}, // Jangan tampilkan untracked files
	}

	for _, cmd := range commands {
		exec.Command(cmd[0], cmd[1:]...).Run()
	}
}

func buatFakeCommit(targetCommits int) {
	fmt.Printf("\033[33mTarget commit: %d dengan kecepatan ultra-high\033[0m\n", targetCommits)

	// Pastikan file test.txt ada
	if _, err := os.Stat("test.txt"); os.IsNotExist(err) {
		file, err := os.Create("test.txt")
		if err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err)
			return
		}
		fmt.Fprintf(file, "Commit for %s\n", time.Now().Format("2006-01-02 15:04:05"))
		file.Close()
	}

	// Counter atomic untuk tracking jumlah commit yang sudah selesai
	var commitsSelesai int64 = 0

	// Channel untuk komunikasi antar goroutine
	batchJobs := make(chan int, targetCommits/BATCH_SIZE+1)

	// WaitGroup untuk menunggu semua goroutine selesai
	var wg sync.WaitGroup

	// Waktu mulai untuk menghitung rate
	startTime := time.Now()

	// Buat worker goroutine
	for w := 1; w <= CONCURRENT_WORKERS; w++ {
		wg.Add(1)
		go batchWorker(w, batchJobs, &commitsSelesai, &wg, targetCommits)
	}

	// Hitung jumlah batch yang dibutuhkan
	numBatches := (targetCommits + BATCH_SIZE - 1) / BATCH_SIZE

	// Kirim jobs batch ke channel
	for j := 0; j < numBatches; j++ {
		startID := j*BATCH_SIZE + 1
		batchJobs <- startID
	}
	close(batchJobs)

	// Goroutine untuk update progress bar
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				count := atomic.LoadInt64(&commitsSelesai)
				if count >= int64(targetCommits) {
					return
				}

				// Hitung kecepatan commit/detik
				elapsedSeconds := time.Since(startTime).Seconds()
				if elapsedSeconds <= 0 {
					elapsedSeconds = 0.001 // Hindari pembagian dengan nol
				}
				rate := float64(count) / elapsedSeconds

				// Tampilan progress bar
				progress := int((float64(count) / float64(targetCommits)) * 30)
				progressBar := strings.Repeat("#", progress) + strings.Repeat(" ", 30-progress)
				percentage := int((float64(count) / float64(targetCommits)) * 100)

				fmt.Printf("\r\033[32m[%s] %d/%d commits \033[36m(%d%%) - %.2f commits/s\033[0m",
					progressBar, count, targetCommits, percentage, rate)

				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	// Tunggu semua worker selesai
	wg.Wait()
	done <- true

	// Hitung total waktu dan rate
	duration := time.Since(startTime)
	rate := float64(targetCommits) / duration.Seconds()

	fmt.Print("\n\n")
	fmt.Printf("\033[32m✓ \033[1mSelesai! %d commits berhasil dibuat dalam %.2f detik (%.2f commits/s)\033[0m\n",
		targetCommits, duration.Seconds(), rate)
	fmt.Print("\033[33mJangan lupa untuk melakukan git push untuk mengirim commit ke repository.\033[0m\n")
	fmt.Print("\033[36mgit push origin <nama-branch>\033[0m\n")
}

func batchWorker(id int, batchJobs <-chan int, commitsSelesai *int64, wg *sync.WaitGroup, targetCommits int) {
	defer wg.Done()

	for startID := range batchJobs {
		// Hitung jumlah commit dalam batch ini
		endID := startID + BATCH_SIZE - 1
		if endID > targetCommits {
			endID = targetCommits
		}

		batchSize := endID - startID + 1

		// Gunakan teknik direct filesystem untuk kecepatan maksimal
		// alih-alih append satu per satu, kita generate satu file besar
		buffer := new(strings.Builder)

		// Header untuk batch ini
		waktuSekarang := time.Now()
		formatWaktu := waktuSekarang.Format("2006-01-02 15:04:05")

		// Generate semua konten sekaligus
		for i := 0; i < batchSize; i++ {
			buffer.WriteString(fmt.Sprintf("Commit %d - %s\n", startID+i, formatWaktu))
		}

		// Tulis semua perubahan ke file dalam satu operasi
		err := os.WriteFile("test.txt", []byte(buffer.String()), 0644)
		if err != nil {
			continue
		}

		// Jalankan perintah git untuk semua perubahan dalam batch
		// Gunakan teknik untuk mempercepat git
		cmdAdd := exec.Command("git", "add", "-A")
		err = cmdAdd.Run()

		if err == nil {
			// Buat commit message yang menggabungkan semua
			commitMessage := fmt.Sprintf("Turbo batch commit %d-%d - %s",
				startID, endID, formatWaktu)

			// Gunakan flag untuk mempercepat commit
			cmdCommit := exec.Command("git", "commit", "-m", commitMessage, "--no-verify", "--no-gpg-sign")
			err = cmdCommit.Run()

			if err == nil {
				// Update counter atomic
				atomic.AddInt64(commitsSelesai, int64(batchSize))
			}
		}
	}
}
