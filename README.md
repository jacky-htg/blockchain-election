# Sistem Pemungutan Suara Berbasis Blockchain

## Gambaran Umum
Proyek ini adalah sistem pemungutan suara berbasis blockchain yang dirancang untuk memastikan transparansi, keamanan, dan efisiensi dalam proses pemilihan. Sistem ini memanfaatkan teknologi terdesentralisasi untuk membuat catatan suara yang tidak dapat diubah dan untuk memfasilitasi proses pemungutan suara. Sistem ini mencakup fitur seperti pendaftaran pemilih, pendaftaran kandidat, dan penanganan suara, semua diimplementasikan dalam jaringan peer-to-peer (P2P).

## Fitur
- Pendaftaran Kandidat: Memfasilitasi kandidat untuk mendaftar dalam pemilihan.
- Proses Pemungutan Suara: Mendukung pemungutan suara yang aman dengan menjaga kerahasiaan pemilih.
- Penghitungan Suara: Secara otomatis menghitung suara dan memberikan hasil.
- Teknologi Blockchain: Menjamin integritas dan ketidakberubahan catatan pemungutan suara.
- Jaringan P2P: Memfasilitasi komunikasi langsung antar node untuk memastikan keandalan dan toleransi kesalahan.

## Teknologi yang Digunakan
- Golang: Bahasa pemrograman utama untuk membangun aplikasi ini.
- Blockchain: Teknologi inti untuk menyimpan dan mengelola suara dengan aman.

## Memulai
Prasyarat
- Instal Go (versi 1.23.1 atau lebih baru)

## Instalasi

**Klon repositori:**

```bash
git clone git@github.com:jacky-htg/blockchain-election.git
cd blockchain-election
```

**Instal dependensi:**

```bash
go mod tidy
```

**Jalankan aplikasi:**

Untuk mnejalankan node pertama (yang akan membuat genesis block)
```bash
go run main.go -address=localhost:3000 -init=true  
```

Untuk menjalankan node baru, buka terminal baru dan jalankan perintah:
```bash
go run main.go -address=localhost:3001
```

Anda bisa membuat peer sebanyak yang anda inginkan dengan address yang berbeda-beda.


## Penggunaan

- **Pemungutan Suara:** Dalam terminal, ketik perintah `vote voterID kandidatID`, misal: `vote user111 Bob`
- **Lihat Hasil:** Dalam terminal, ketik perintah `showresult`


## Kontribusi
Kontribusi sangat diterima! Jika Anda memiliki saran atau perbaikan, silakan ajukan pull request atau buka isu.

## Lisensi
Proyek ini dilisensikan di bawah Lisensi GNU GPL - lihat berkas [LICENSE](./LICENSE) untuk detailnya.

## Kontak
Untuk pertanyaan atau pertanyaan lebih lanjut, silakan hubungi [rijal.asep.nugroho@gmail.com].