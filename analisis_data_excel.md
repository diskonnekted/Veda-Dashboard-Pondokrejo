# Laporan Analisis Data Excel
## File: `1 KK_ART Pondokrejo.xlsx`
## Potensi Data yang Dapat Ditampilkan di VEDA Dashboard

> **Statistik File:**
> - **Total Sheet:** 1 (Laporan Keluarga)
> - **Total Kolom:** 295 kolom
> - **Total Baris Data:** ± 5.900 baris (KK + anggota/ART)
> - **Format Header:** 3 baris (Row 1: judul, Row 2: kode variabel, Row 3: nama kolom)
> - **Data dimulai:** Row 4

---

## KELOMPOK A — IDENTITAS & ADMINISTRASI RUMAH TANGGA

| Kolom | Keterangan | Potensi Tampilan |
|---|---|---|
| NO KK | Nomor Kartu Keluarga | Pencarian & filter KK |
| ID Ruta | ID Rumah Tangga internal | Key identifier |
| Nama Kepala KK | Nama kepala keluarga | Marker popup peta, tabel data |
| Alamat / Alamat Lengkap | Alamat domisili | Detail popup |
| RW / RT | Nomor RT/RW | Filter per RT/RW di peta |
| ID Dusun | Kode Padukuhan | Filter per Dusun |
| Koordinat | Lat,Lng GPS | 🗺️ TITIK MARKER PETA |
| Keterangan | Catatan khusus | Popup detail |
| Tgl Pendataan | Tanggal survei | Filter periode, validasi aktualitas |
| Tgl Pemeriksaan | Tanggal verifikasi | Audit trail |

---

## KELOMPOK B — STATUS KESEJAHTERAAN & KATEGORI

| Kolom | Keterangan | Potensi Tampilan |
|---|---|---|
| ID Miskin | Status miskin (1/2) | 🎯 Filter kemiskinan |
| ID Desil | Desil 1-4 kesejahteraan | 📊 Warna marker peta, grafik desil |
| Percentile | Persentil skor kemiskinan | Grafik distribusi |
| Mampu | Flag keluarga mampu | Segmentasi peta |
| IS Usulan | Diusulkan sebagai penerima baru | Daftar prioritas usulan baru |
| IS Dtks | Terdaftar di DTKS | Validasi cross-check DTKS |
| IS Ekstrem | Kategori Miskin Ekstrem | ⚠️ Filter prioritas absolut |
| IS Stun | Keluarga berisiko Stunting | 🍼 Dashboard Stunting |
| IS Verify | Status verifikasi lapangan | Badge verifikasi di tabel |
| IS Musdus | Hasil musyawarah dusun | Transparansi musyawarah |
| IS Musdes | Hasil musyawarah desa | Transparansi musyawarah |

---

## KELOMPOK C — KONDISI PERUMAHAN & INFRASTRUKTUR

| Kolom | Keterangan | Potensi Tampilan |
|---|---|---|
| Lantai Luas | Luas lantai rumah (m²) | Grafik distribusi luas hunian |
| ID Lantai | Jenis lantai | Filter RTLH, grafik donat |
| ID Dinding | Jenis dinding | Filter RTLH |
| ID Atap | Jenis atap | Filter RTLH |
| ID Airminum | Sumber air minum | 💧 Peta akses air bersih |
| ID Airminum Jarak | Jarak ke sumber air | Analisis aksesibilitas |
| ID Listrik | Akses listrik | ⚡ Peta elektrifikasi |
| ID Listrik Daya | Daya listrik (watt) | Grafik distribusi daya |
| ID Fasbab | Fasilitas BAB | 🚽 Dashboard sanitasi |
| ID Kloset | Jenis kloset | Detail sanitasi |
| ID Tinja | Pembuangan tinja | Detail sanitasi |
| Sampah A-E | Pengelolaan sampah | Dashboard lingkungan |
| ID Bbm | Bahan bakar masak | Grafik jenis BBM |
| Rumah Sutet | Dekat SUTET | Peta risiko lingkungan |
| Rumah Sei | Dekat Sungai | Peta risiko banjir |
| Polusi Air/Tanah/Udara | Status pencemaran | ⚠️ Peta risiko lingkungan |

---

## KELOMPOK D — BANTUAN SOSIAL YANG DITERIMA

| Kolom | Keterangan | Potensi Tampilan |
|---|---|---|
| IS Bpnt / Bpnt Thn | Penerima BPNT (Sembako) | 🛒 Dashboard BPNT |
| IS Pkh / Pkh Thn | Penerima PKH | 👨‍👩‍👧 Dashboard PKH |
| IS Blt / Blt Thn | Penerima BLT | Dashboard BLT |
| IS Listrik | Subsidi listrik | Dashboard Subsidi |
| IS Banpem | Bantuan Pemerintah lainnya | Rekap Banpem |
| IS Pupuk | Subsidi pupuk | Dashboard pertanian |
| IS Lpg | Subsidi LPG | Dashboard energi |
| IS Baznas | Bantuan BAZNAS | Dashboard Sosial |
| IS Csr | Bantuan CSR | Dashboard CSR |

---

## KELOMPOK E — ASET & PEREKONOMIAN

| Kolom | Keterangan | Potensi Tampilan |
|---|---|---|
| Atis Rumah | Status kepemilikan rumah | Grafik RTLH & kepemilikan |
| AT Tani Sawah | Luas lahan sawah | 🌾 Dashboard Pertanian |
| AT Tani Nonsawah | Lahan non-sawah | Analisis pertanian |
| AB Gas/TV/Laptop/Motor/Mobil/HP | Kepemilikan aset | Grafik aset rumah tangga |
| Ekor Sapi/Kerbau/Kambing | Ternak | 🐄 Dashboard Peternakan |
| IS Inet | Akses internet | 📡 Peta Konektivitas Digital |
| IS Bank | Rekening bank | Inklusi keuangan |
| Overall Sum | Total pengeluaran | 📊 Distribusi pengeluaran |
| Makan Avg / Brg Avg | Rata-rata pengeluaran | Analisis kemampuan ekonomi |
| Uang Rokok | Pengeluaran rokok | Analisis konsumsi |

---

## KELOMPOK F — USAHA MIKRO & KETENAGAKERJAAN

| Kolom | Keterangan | Potensi Tampilan |
|---|---|---|
| IS Ush | Punya usaha | 🏪 Dashboard UMKM |
| Ush Detail | Jenis usaha | Profil UMKM per dusun |
| Ush Omset | Omset usaha | Analisis ekonomi mikro |
| Ush Izin | Usaha punya izin/NIB | Dashboard legalitas usaha |
| Ush Inet | Usaha berbasis internet | Indeks UMKM Digital |
| Mikro Tekstil/Kayu/Makan... | Jenis industri mikro | Peta industri per wilayah |
| Sarana Market/Resto/Toko/Bengkel... | Fasilitas ekonomi | 🗺️ Peta fasilitas ekonomi |

---

## KELOMPOK G — DATA ANGGOTA RUMAH TANGGA (ART / Per Individu)

| Kolom | Keterangan | Potensi Tampilan |
|---|---|---|
| NIK | Nomor Induk Kependudukan | Pencarian warga |
| Nama | Nama lengkap | Detail popup |
| ID Hub Keluarga | Hubungan dengan kepala KK | Filter kepala/istri/anak |
| Tgl Lahir / Usia | Tanggal lahir & usia | 📅 Piramida usia penduduk |
| ID Kelamin | Jenis kelamin | Grafik gender |
| ID Agama | Agama | Statistik keagamaan |
| ID Nikah | Status pernikahan | Grafik pernikahan |
| Hamil / Tgl Hpl | Status kehamilan | 🤰 Dashboard ibu hamil |
| ID Jenjang / ID Ijazah | Pendidikan | 🏫 Dashboard pendidikan |
| ID Penyakit Kronis | Penyakit kronis | 🏥 Dashboard kesehatan |
| ID Difable | Status disabilitas | Dashboard disabilitas |
| ID Kerja / Kerja Detail | Pekerjaan | 💼 Dashboard ketenagakerjaan |
| Income | Penghasilan individu | Distribusi pendapatan |
| IS Rokok | Perokok | Analisis konsumsi rokok |
| ID Gizi | Status gizi | Stunting monitoring |
| ID Asi | Status ASI eksklusif | Dashboard kesehatan ibu-anak |
| Kes Lihat/Dengar/Jalan/Tangan... | Kesulitan fungsional | Dashboard disabilitas detail |
| Sos Jamkes | Jaminan kesehatan | Dashboard BPJS |
| Sos Prakerja | Pernah ikut Prakerja | Dashboard ketenagakerjaan |
| Sos Kur | Penerima KUR | Dashboard keuangan inklusif |
| Sos Pip | Penerima PIP | Dashboard beasiswa |
| Sos Jamket | Jaminan ketenagakerjaan | Dashboard perlindungan kerja |

---

## RINGKASAN: 10 FITUR BARU YANG DIREKOMENDASIKAN

| # | Fitur Baru | Sumber Kolom |
|---|---|---|
| 1 | **Dashboard DTKS & Validasi** | IS Dtks, IS Verify, IS Musdus, IS Musdes |
| 2 | **Peta Konektivitas Digital** | IS Inet, IS Bank, AB Laptop, AB HP, Ush Inet |
| 3 | **Dashboard Stunting & Ibu Hamil** | IS Stun, Hamil, ID Gizi, ID Asi, Usia < 5 |
| 4 | **Peta Risiko Lingkungan** | Rumah Sei, Polusi Air/Tanah/Udara |
| 5 | **Dashboard UMKM Lengkap** | IS Ush, Ush Detail, Ush Omset, Ush Izin |
| 6 | **Piramida Penduduk** | Usia, ID Kelamin per KK |
| 7 | **Dashboard Perlindungan Sosial** | IS PKH, BPNT, BLT, PIP, Jamkes, Prakerja |
| 8 | **Indeks Kemiskinan Multi-Dimensi** | Gabungan Desil + Infrastruktur + Aset |
| 9 | **Dashboard Pertanian & Peternakan** | AT Tani Sawah, Ekor Sapi/Kambing |
| 10 | **Dashboard Disabilitas** | ID Difable, Kes Lihat/Dengar/Jalan, dll |

---

*Laporan ini dibuat berdasarkan analisis programatik file `1 KK_ART Pondokrejo.xlsx` (295 kolom, ±5.900 baris).*
*Dibuat: 10 Mei 2026 — Antigravity AI Coding Assistant*
