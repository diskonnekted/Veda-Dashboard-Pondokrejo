# Laporan Teknis Aplikasi: Veda Dashboard Pondokrejo

## 1. Pendahuluan
**Veda Dashboard Pondokrejo** (Visual Economic Data Analytics) adalah aplikasi berbasis web yang dirancang untuk visualisasi, analisis, dan manajemen data kependudukan serta ekonomi di tingkat desa (Pondokrejo). Aplikasi ini mengintegrasikan data tabular (Excel) dengan Sistem Informasi Geografis (GIS) untuk memberikan wawasan yang mendalam bagi pengambil kebijakan desa.

---

## 2. Arsitektur Sistem

### 2.1 Backend (Go & Gin Framework)
Aplikasi dibangun menggunakan bahasa pemrograman **Go** dengan framework **Gin-Gonic**. Pemilihan teknologi ini didasarkan pada performa tinggi, efisiensi memori, dan kemudahan dalam pengembangan API.

- **Entry Point**: `main.go` mengelola routing, konfigurasi server (default port 8080), dan penyajian aset statis.
- **Data Handling**: Aplikasi tidak menggunakan database SQL tradisional, melainkan mengandalkan file Excel sebagai sumber data utama, yang diproses secara dinamis saat aplikasi dijalankan.

### 2.2 Frontend (HTML, Tailwind CSS, JS)
Antarmuka pengguna menggunakan template HTML yang disajikan secara server-side oleh Gin.
- **Styling**: Menggunakan **Tailwind CSS** untuk desain modern dan responsif.
- **Visualisasi Data**: Menggunakan **Chart.js** untuk grafik statistik dan **Leaflet.js** untuk peta interaktif.
- **Aestetik**: Mengadopsi prinsip desain premium dengan transparansi (glassmorphism) dan palet warna yang harmonis.

---

## 3. Komponen Utama & Logika Teknis

### 3.1 Pengolahan Data (Parser)
File `parser.go` bertanggung jawab untuk:
- **Excel Parsing**: Membaca file `KK_Data_Final_Readable.xlsx` menggunakan library `excelize`.
- **Struktur Data**: Memetakan baris Excel ke dalam struct `Household` (Rumah Tangga) dan `Resident` (Penduduk).
- **Normalisasi**: Membersihkan data mentah, menangani konversi tipe data, dan memetakan kode RW ke nama Padukuhan (Dusun).

### 3.2 Integrasi GIS & Geospasial
Salah satu fitur tercanggih adalah validasi koordinat otomatis di `parser.go`:
- **Point-in-Polygon**: Menggunakan algoritma *Ray Casting* untuk memverifikasi apakah koordinat rumah tangga berada di dalam batas wilayah Desa Pondokrejo (`PONDOKREJO.geojson`).
- **Auto-Correction**: Jika koordinat terdeteksi di luar batas desa (akibat kesalahan input GPS), aplikasi secara otomatis menggeser titik tersebut ke titik pusat (centroid) Dusun yang bersangkutan.

### 3.3 Dashboard Analitik
File `analytics.go` melakukan agregasi data secara real-time untuk menghasilkan metrik:
- **Kesejahteraan**: Distribusi pendapatan dan profil pekerjaan penduduk miskin (Desil 1 & 2).
- **Infrastruktur**: Statistik Rumah Tidak Layak Huni (RTLH), kepemilikan jamban, dan akses air bersih.
- **Kelompok Rentan**: Identifikasi Lansia tunggal dan keluarga miskin yang memiliki balita.

### 3.4 Editor Layer Geospasial
Aplikasi dilengkapi dengan fitur `editor` yang memungkinkan admin untuk mengelola layer peta (seperti jalan, sawah, bangunan) yang disimpan dalam format GeoJSON di folder `/layers`.

---

## 4. Keamanan & Akses
- **Otentikasi**: Menggunakan sistem login sederhana berbasis sesi (`sessionStorage`) untuk mengamankan akses ke dashboard dan editor.
- **Environment**: Konfigurasi port dapat disesuaikan melalui environment variable `PORT`.

---

## 5. Alur Deployment (Build Process)
Aplikasi mendukung mode eksekusi langsung (`go run .`) atau dikompilasi menjadi binary executable tunggal (`veda.exe`).
- **Static Generation**: Terdapat flag `-gen` untuk menghasilkan file JSON statis (`residents.json`) yang dapat digunakan untuk deployment pada host statis jika backend Go tidak digunakan secara permanen.

---

## 6. Daftar File Penting
- `main.go`: Inti aplikasi dan routing.
- `parser.go`: Logika ekstraksi data Excel dan koreksi GIS.
- `analytics.go`: Logika kalkulasi metrik statistik.
- `templates/`: Folder berisi halaman UI (Login, Dashboard, Analytics, Editor).
- `layers/`: Folder penyimpanan data geospasial tematik.

---
**Dibuat oleh:** Antigravity AI Coding Assistant
**Tanggal:** 10 Mei 2026
