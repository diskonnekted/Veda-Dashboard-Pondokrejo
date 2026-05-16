# LAMPIRAN I
# SPESIFIKASI TEKNIS SISTEM
## VEDA Dashboard — Visual Economic Data Analytics
### Untuk Dinas Sosial Kabupaten Banjarnegara

---

## A. Identifikasi Sistem

| Atribut | Keterangan |
|---|---|
| **Nama Sistem** | VEDA Dashboard (Visual Economic Data Analytics) |
| **Versi** | 1.2.0 |
| **Tipe Aplikasi** | Web Application (Server-side Rendered) |
| **Bahasa Pemrograman** | Go (Golang) v1.25 |
| **Framework Backend** | Gin-Gonic v1.12 |
| **Format Data Input** | Microsoft Excel (.xlsx), GeoJSON (.geojson) |
| **Format Data Output** | JSON, GeoJSON, HTML (Laporan Cetak) |
| **Mode Operasi** | Lokal (Offline) / Jaringan Intranet / Internet |

---

## B. Spesifikasi Perangkat Keras (Hardware Requirements)

### B.1 Server / Komputer Operator (Minimum)
| Komponen | Spesifikasi Minimum | Spesifikasi Rekomendasi |
|---|---|---|
| **Prosesor** | Intel Core i3 Gen 8 / AMD Ryzen 3 | Intel Core i5 Gen 10 / AMD Ryzen 5 |
| **RAM** | 4 GB DDR4 | 8 GB DDR4 |
| **Penyimpanan** | 50 GB HDD | 100 GB SSD |
| **OS** | Windows 10 64-bit / Ubuntu 20.04 | Windows 11 / Ubuntu 22.04 LTS |
| **Koneksi Jaringan** | Opsional (bisa offline) | 10 Mbps (untuk deployment multi-user) |

### B.2 Klien / Perangkat Akses (Browser)
| Browser | Versi Minimum |
|---|---|
| Google Chrome | 90+ |
| Mozilla Firefox | 88+ |
| Microsoft Edge | 90+ |
| Safari | 14+ |
| **Resolusi Layar** | 1280 x 768 (minimum) |

---

## C. Arsitektur Sistem

```
┌─────────────────────────────────────────────┐
│              PENGGUNA (Browser)             │
│    Dashboard │ Peta │ Analitik │ Editor    │
└──────────────────┬──────────────────────────┘
                   │ HTTP/HTTPS
┌──────────────────▼──────────────────────────┐
│         WEB SERVER — Gin-Gonic (Go)         │
│  ┌─────────────┐   ┌──────────────────────┐ │
│  │   Router    │   │   Template Engine    │ │
│  │  (main.go)  │   │   (HTML Templates)   │ │
│  └──────┬──────┘   └──────────────────────┘ │
│         │                                   │
│  ┌──────▼──────┐   ┌──────────────────────┐ │
│  │  Data Layer │   │  GIS / Geo Engine    │ │
│  │ (parser.go) │   │   (geo_clip.go)      │ │
│  └──────┬──────┘   └──────────────────────┘ │
│         │                                   │
│  ┌──────▼──────┐   ┌──────────────────────┐ │
│  │  Analytics  │   │  Editor State Store  │ │
│  │(analytics.go│   │ (editor_storage.go)  │ │
│  └─────────────┘   └──────────────────────┘ │
└─────────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│            SUMBER DATA                      │
│  📊 Excel (.xlsx)  🗺️ GeoJSON  📁 Layers/   │
└─────────────────────────────────────────────┘
```

---

## D. Modul Fungsional Sistem

### D.1 Modul Pemetaan & GIS

| Fitur | Deskripsi Teknis |
|---|---|
| **Peta Interaktif** | Menggunakan Leaflet.js dengan layer tile OpenStreetMap |
| **Sebaran KK** | Setiap KK ditampilkan sebagai marker berdasarkan koordinat lat/lng |
| **Kode Warna Desil** | Marker diberi warna berbeda: Merah (Miskin Ekstrem), Oranye (Miskin), Kuning (Hampir Miskin), Hijau (Mampu) |
| **Filter Dusun/RW** | Pengguna dapat menyaring data per padukuhan/RW secara real-time |
| **Klik Marker** | Menampilkan popup detail KK: Nama KK, Anggota, Bantuan, Kondisi Rumah |
| **Validasi Batas Wilayah** | Algoritma *Ray Casting Point-in-Polygon* berbasis GeoJSON batas desa |
| **Auto-Correction GPS** | Koordinat di luar batas desa digeser ke centroid dusun secara otomatis |
| **Layer Tematik** | Jalan, Sungai, Sawah, Irigasi, Bangunan, TPS — dapat diaktifkan/nonaktifkan |

### D.2 Modul Analitik Sosial-Ekonomi

| Metrik | Cara Penghitungan |
|---|---|
| **Distribusi Pendapatan** | Kategorisasi pengeluaran: < Rp1 Juta / Rp1-2 Juta / > Rp2 Juta |
| **Profil Pekerjaan Miskin** | Pekerjaan kepala KK pada Desil 1 & 2 dikelompokkan dan dihitung frekuensinya |
| **Tingkat Pendidikan Miskin** | Pendidikan kepala KK pada Desil 1 & 2 |
| **RTLH** | KK dengan status kepemilikan bukan "Milik Sendiri" |
| **Tanpa Jamban** | KK dengan kolom sanitasi bukan "Milik Sendiri" |
| **Tanpa Air Bersih** | KK yang sumber airnya mengandung kata "Sungai" atau "Danau" |
| **Lansia Tunggal** | Kepala KK usia >65 tahun dengan anggota ≤ 2 jiwa |
| **Keluarga Miskin + Balita** | KK Desil 1-2 yang memiliki anggota berusia < 5 tahun |

### D.3 Modul Editor Geospasial

| Fitur | Keterangan |
|---|---|
| **CRUD Layer** | Tambah, ubah, hapus layer GeoJSON secara visual dari browser |
| **Persistensi State** | Perubahan disimpan ke file `editor_state.json` di server |
| **Preview Real-time** | Setiap perubahan langsung tercermin di peta tanpa refresh |

### D.4 Modul Import Data

| Aspek | Detail |
|---|---|
| **Format Diterima** | Microsoft Excel (.xlsx) — format DTKS/PKH standar |
| **Library Pengolah** | `excelize/v2` — library Go berperforma tinggi |
| **Kapasitas** | Diuji hingga > 150 kolom dan > 10.000 baris tanpa degradasi performa signifikan |
| **Normalisasi** | Pembersihan otomatis data kosong, spasi ganda, dan awalan apostrof dari Excel |

---

## E. Keamanan Sistem

| Aspek | Implementasi |
|---|---|
| **Autentikasi** | Login berbasis sesi browser (`sessionStorage`) dengan validasi kredensial di server |
| **Akses Halaman** | Setiap halaman sensitif diperiksa status login sebelum dimuat |
| **Isolasi Data** | Data berjalan di server lokal, tidak terekspos ke internet kecuali dikonfigurasi |
| **Upgrade Keamanan** | Tersedia opsi implementasi JWT-based authentication & HTTPS untuk deployment publik |

---

## F. Kemampuan Deployment

| Mode | Keterangan |
|---|---|
| **Lokal (Offline)** | Berjalan di satu komputer tanpa jaringan; cocok untuk operator tunggal |
| **Intranet LAN** | Deploy di jaringan kantor; dapat diakses oleh seluruh staf Dinsos dari browser masing-masing |
| **Cloud/VPS** | Dapat di-deploy ke server Linux untuk akses online kabupaten-wide |
| **Static Export** | Flag `--gen` menghasilkan file JSON statis yang dapat di-host di web hosting biasa |
| **Executable Tunggal** | Dikompilasi menjadi satu file `.exe` (Windows) atau binary (Linux); tidak perlu install runtime tambahan |

---

## G. Integrasi & Kompatibilitas

| Sistem Eksternal | Status Kompatibilitas |
|---|---|
| **DTKS Kemensos** | ✅ Kompatibel melalui ekspor Excel |
| **Aplikasi PKH** | ✅ Data PKH dapat diimpor via format Excel |
| **SIKS-NG** | ✅ Format ekspor SIKS-NG (.xlsx) dapat diolah langsung |
| **OpenStreetMap** | ✅ Digunakan sebagai basemap peta tanpa biaya lisensi |
| **BIG (Badan Informasi Geospasial)** | ✅ Mendukung format GeoJSON standar nasional |

---

## H. Service Level Agreement (SLA) yang Ditawarkan

| Layanan | Cakupan |
|---|---|
| **Instalasi & Konfigurasi** | Pemasangan aplikasi di server/komputer Dinsos |
| **Migrasi Data** | Konversi dan pembersihan data Excel existing ke format sistem |
| **Pelatihan Pengguna** | 2 sesi pelatihan (operator data & kepala bidang) |
| **Garansi Bug Fix** | 3 bulan pasca instalasi |
| **Pemeliharaan Tahunan** | Pembaruan fitur, penyesuaian kolom data, dan dukungan teknis |
| **Response Time** | Maksimal 2×24 jam untuk laporan gangguan kritis |

---

## I. Daftar Dependensi (Teknologi yang Digunakan)

| Komponen | Teknologi | Lisensi |
|---|---|---|
| Backend Runtime | Go 1.25 | BSD 3-Clause (Free) |
| Web Framework | Gin-Gonic v1.12 | MIT License (Free) |
| Excel Parser | excelize/v2 | BSD License (Free) |
| Peta Interaktif | Leaflet.js | BSD 2-Clause (Free) |
| Grafik Statistik | Chart.js | MIT License (Free) |
| UI Styling | Tailwind CSS CDN | MIT License (Free) |
| Basemap Peta | OpenStreetMap | ODbL (Free) |

> **Catatan:** Seluruh teknologi yang digunakan berstatus **Free & Open Source**, sehingga tidak ada biaya lisensi pihak ketiga yang dibebankan kepada Dinas Sosial Kabupaten Banjarnegara.

---

**Lampiran ini merupakan bagian tidak terpisahkan dari Surat Penawaran VEDA Dashboard.**

**Banjarnegara, 10 Mei 2026**

**Tim Pengembangan VEDA Dashboard — Clasnet Group**
