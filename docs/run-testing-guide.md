# FELO Testing Guide

## Tujuan

Dokumen ini menjelaskan cara menjalankan unit test FELO secara lokal dan di Jenkins.
Scope saat ini mencakup scaffold test-first untuk:

- `ride-service`
- `matching-service`
- `wallet-service`
- `payment-service`
- `location-service`

## Struktur Repo yang Relevan

```text
services/                     Source service dan unit test
tools/coveragecheck/          Validasi threshold coverage
tools/gotest2junit/           Konversi output go test ke JUnit XML
docs/testing-strategy.md      Standar testing
Jenkinsfile                   Pipeline Jenkins
```

## Prasyarat

Pastikan hal berikut tersedia:

1. Go sudah terpasang.
2. Command `go` bisa dijalankan dari terminal.
3. Working directory berada di root repo

## 1. Menjalankan Semua Unit Test

Gunakan perintah ini:

```powershell
go test ./...
```

Hasil yang diharapkan:

- semua package `internal/service` berstatus `ok`
- package `domain` dan `ports` bisa tampil sebagai `[no test files]`

## 2. Menjalankan Unit Test per Service

Contoh per service:

```powershell
go test ./services/ride-service/...
go test ./services/matching-service/...
go test ./services/wallet-service/...
go test ./services/payment-service/...
go test ./services/location-service/...
```

Gunakan ini saat Anda hanya sedang fokus pada satu service.

## 3. Menjalankan Test Verbose

Kalau ingin melihat nama test case satu per satu:

```powershell
go test -v ./services/ride-service/...
```

## 4. Menjalankan Satu Test Case Tertentu

Contoh:

```powershell
go test -v ./services/ride-service/... -run TestTripService_StartRide_ArrivedTrip_PublishesRideStartedEvent
```

Ini berguna saat sedang debugging satu behavior saja.

## 5. Menjalankan Coverage Report

Untuk generate coverage file:

```powershell
go test -covermode=atomic -coverprofile='coverage.out' ./services/...
```

Untuk melihat coverage per function:

```powershell
go tool cover -func 'coverage.out'
```

Untuk generate HTML report:

```powershell
go tool cover -html='coverage.out' -o 'coverage.html'
```

Setelah itu buka file `coverage.html` di browser.

## 6. Menjalankan Threshold Coverage Check

Repo ini punya helper internal untuk memastikan coverage tidak di bawah batas minimum.

Contoh:

```powershell
go run ./tools/coveragecheck -file 'coverage.out' -threshold 70
```

Arti threshold saat ini:

- `70` untuk batas minimum umum
- target ideal business logic tetap `80` atau lebih

Kalau coverage di bawah threshold, command akan gagal dan cocok dipakai sebagai gate di Jenkins.

## 7. Generate JUnit XML untuk Jenkins

Jenkins biasanya lebih mudah membaca format `JUnit XML`.
Gunakan langkah berikut:

```powershell
go test -json ./services/... | Tee-Object -FilePath 'gotest.json'
Get-Content -LiteralPath 'gotest.json' | go run ./tools/gotest2junit | Set-Content 'junit.xml'
```

File hasil:

- `gotest.json`
- `junit.xml`

File `junit.xml` inilah yang nanti dibaca oleh Jenkins.

## 8. Menjalankan Alur Lengkap Lokal

Ini adalah urutan yang paling mendekati Jenkins:

```powershell
go test -json -covermode=atomic -coverprofile='coverage.out' ./services/... | Tee-Object -FilePath 'gotest.json'
go test ./tools/...
go tool cover -html='coverage.out' -o 'coverage.html'
Get-Content -LiteralPath 'gotest.json' | go run ./tools/gotest2junit | Set-Content 'junit.xml'
go run ./tools/coveragecheck -file 'coverage.out' -threshold 70
```

Kalau semua sukses, berarti lokal Anda sudah siap untuk CI.

## 9. Menjalankan Race Check

Untuk environment yang mendukung:

```powershell
go test -race ./...
```

Catatan penting:

Environment lokal saat ini terdeteksi sebagai `windows/386`, dan kombinasi ini tidak mendukung `-race`.
Karena itu Jenkinsfile sudah dibuat supaya `-race` otomatis dilewati pada environment seperti ini.

Kalau nanti Jenkins agent Anda memakai environment seperti berikut, `-race` sebaiknya diaktifkan:

- `linux/amd64`
- `windows/amd64`
- `darwin/amd64`
- `darwin/arm64`

## 10. Menjalankan dari Jenkins

Pipeline dasar sudah ada di [Jenkinsfile](C:\Users\Harri Supriadi\Documents\unit-test-felo\Jenkinsfile).

Yang dilakukan pipeline:

1. cek toolchain Go
2. jalankan unit test service
3. jalankan test untuk helper tool internal
4. generate coverage report
5. generate `junit.xml`
6. fail build jika coverage di bawah threshold
7. archive artifact hasil test

Artifact utama yang dihasilkan:

- `coverage.out`
- `coverage.html`
- `gotest.json`
- `junit.xml`

## 11. Workflow Harian yang Disarankan

Saat implementasi service nanti, alur kerja yang saya sarankan:

1. tulis atau update unit test dulu
2. jalankan test untuk service terkait saja
3. perbaiki implementasi sampai test pass
4. jalankan `go test ./...`
5. jalankan coverage check
6. sebelum push, jalankan alur lengkap lokal

## 12. Troubleshooting

### Error `-race is not supported on windows/386`

Penyebab:

- arsitektur Go lokal tidak mendukung race detector

Solusi:

- jalankan test tanpa `-race` di lokal
- gunakan Jenkins agent `amd64` untuk race check

### Error coverage file tidak ditemukan

Penyebab umum:

- command `go test -coverprofile='coverage.out'` belum dijalankan
- working directory bukan root repo

Solusi:

```powershell
cd 'C:\Users\Harri Supriadi\Documents\unit-test-felo'
go test -covermode=atomic -coverprofile='coverage.out' ./services/...
```

### Error `junit.xml` kosong atau tidak terbentuk

Penyebab umum:

- file `gotest.json` tidak terbentuk
- output `go test` tidak dijalankan dengan flag `-json`

Solusi:

```powershell
go test -json ./services/... | Tee-Object -FilePath 'gotest.json'
Get-Content -LiteralPath 'gotest.json' | go run ./tools/gotest2junit | Set-Content 'junit.xml'
```

## 13. File Referensi

- [README.md](C:\Users\Harri Supriadi\Documents\unit-test-felo\README.md)
- [testing-strategy.md](C:\Users\Harri Supriadi\Documents\unit-test-felo\docs\testing-strategy.md)
- [Jenkinsfile](C:\Users\Harri Supriadi\Documents\unit-test-felo\Jenkinsfile)
