# Genius - Proje Iyilestirme Raporu

Bu rapor, projenin kaynak kodunun detayli incelenmesi sonucu tespit edilen iyilestirme noktalarini icerir.

---

## 1. Hata Yonetimi (Error Handling)

### 1.1 `cmd/info.go` - `log.Println` + `panic` Anti-Pattern

`info` komutunun `Run` fonksiyonunda hata yonetimi tutarsiz ve tehlikeli. Hemen hemen her hata durumunda once `log.Println(err)` ardindan `panic(err)` cagriliyor. Bu yaklasim:

- Kullaniciya cirkin bir stack trace gosterir
- Programi ani ve kontrolsuz sekilde sonlandirir
- CLI araclari icin uygun degildir

**Mevcut kod (info.go:46-49):**
```go
hInfo, err := host.Info()
if err != nil {
    log.Println(err)
    panic(err)
}
```

**Onerilen yaklasim:**
```go
hInfo, err := host.Info()
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: failed to get host info: %v\n", err)
    os.Exit(1)
}
```

Veya daha iyisi, `RunE` kullanarak:
```go
RunE: func(cmd *cobra.Command, args []string) error {
    hInfo, err := host.Info()
    if err != nil {
        return fmt.Errorf("failed to get host info: %w", err)
    }
    // ...
}
```

### 1.2 `cmd/cpu.go` - Sessiz Hata Yutma

`cpu.Counts` cagrilarinda hatalar `_` ile tamamen yok sayiliyor:

```go
physicalCores, _ := cpu.Counts(false)
logicalCores, _ := cpu.Counts(true)
```

Bu hatalar en azindan loglanmali veya kullaniciya bildirilmelidir.

### 1.3 `cmd/info.go` - Guvenli Olmayan Type Assertion

Satir 142'de `addr.(*net.IPNet)` type assertion'i kontrol edilmeden kullaniliyor. Eger `addr` farkli bir tipse program panic ile coker:

```go
ipv4 := addr.(*net.IPNet).IP.To4()  // TEHLIKELI
```

**Guvenli hali:**
```go
ipNet, ok := addr.(*net.IPNet)
if !ok {
    continue
}
ipv4 := ipNet.IP.To4()
```

---

## 2. Kod Organizasyonu ve Mimari

### 2.1 `cmd/info.go` - Devasa Monolitik Fonksiyon

`info` komutunun `Run` fonksiyonu yaklasik 160 satir ve tek bir fonksiyon icinde:
- Host bilgisi toplama
- CPU bilgisi toplama
- Bellek bilgisi toplama
- Disk bilgisi toplama
- Ag arayuzu bilgisi toplama
- DNS bilgisi toplama
- NTP bilgisi toplama
- Paket versiyon bilgisi toplama

Her biri ayri bir fonksiyona cikarilmali. Ornegin:

```go
func printHostInfo() error { ... }
func printCPUInfo() error { ... }
func printMemoryInfo() error { ... }
func printDiskInfo() error { ... }
func printNetworkInfo() error { ... }
func printDNSInfo() error { ... }
func printNTPInfo() error { ... }
func printPackageVersions() error { ... }
```

### 2.2 `FolderSize` Tipi Yanlis Pakette

`FolderSize` struct'i `helpers/file.go` icinde tanimlanmis, ancak projede zaten `types/` paketi mevcut. Tum tip tanimlari `types/` altinda toplanmali.

### 2.3 Kullanilmayan Kod

- `GetLargestFolders` fonksiyonu (non-concurrent versiyon) `helpers/file.go`'da mevcut ama hicbir yerde kullanilmiyor. Concurrent versiyonu tercih edilmis. Olume kod kaldirilmali.
- `rootCmd.Flags().BoolP("toggle", "t", false, ...)` - root komutunda tanimli ama kullanilmayan bir flag.
- Bircok dosyada yorum satirlari icinde birakilan eski kod parcalari (cobra scaffold comment'leri).

---

## 3. NTP Yapilandirma Parser'i Hatalari (`helpers/file.go`)

### 3.1 Yorum Satirlarini Kaldirma Hatasi

Mevcut kod, NTP dosyasindaki ilk `#` karakterinden sonrasini tamamen siliyor:
```go
if strings.Index(stringNTPFile, "#") != -1 {
    stringNTPFile = stringNTPFile[:strings.Index(stringNTPFile, "#")]
}
```

Bu yaklasim, eger bir yorum satiri dosyanin basinda server satirlarindan once yer aliyorsa, tum server satirlarini da siler. Dogru yaklasim, dosyayi satir satir okuyup her satirdaki yorumlari ayri ayri kaldirmaktir.

### 3.2 Yalnizca `/etc/ntp.conf` Destegi

Chrony (`/etc/chrony/chrony.conf`) destegi yorum satiri olarak birakilmis. Modern Linux dagitimlarinin cogu chrony kullaniyor. Her ikisi de desteklenmeli:
```go
ntpPaths := []string{"/etc/ntp.conf", "/etc/chrony/chrony.conf", "/etc/chrony.conf"}
```

---

## 4. Hardcoded Ag Arayuzu Isimleri

`cmd/info.go` satir 135'te ag arayuzu isimleri sabit kodlanmis:
```go
if i.Name == "en0" || i.Name == "ens160" {
```

Bu sadece belirli sistemlerde calisir. `eth0`, `wlan0`, `enp0s3`, `ens33`, `wlp2s0` gibi yaygon arayuz isimleri desteklenmiyor. Onerilen yaklasim:

```go
// Loopback olmayan, aktif, IPv4 adresi olan tum arayuzleri goster
for _, i := range ifaces {
    if i.Flags&net.FlagLoopback != 0 || i.Flags&net.FlagUp == 0 {
        continue
    }
    // ...
}
```

---

## 5. Hardcoded NTP Sunucusu

`cmd/info.go` satir 171:
```go
ntpTime, err := ntp.Time("0.tr.pool.ntp.org")
```

Turkiye'ye ozgu bir NTP sunucusu sabit kodlanmis. Bu:
- Farkli ulkelerdeki kullanicilar icin yavas olabilir
- Sunucu erisilemezse program hata verir

Oneriler:
- Flag veya environment variable ile yapilandirilabilir olmali
- Varsayilan olarak `pool.ntp.org` gibi global bir havuz kullanilmali
- Timeout eklenmeli (su an network timeout'a kadar bekliyor)

---

## 6. Test Eksikligi

Projede **hicbir test dosyasi yok**. Bu, en kritik iyilestirme noktalarindan biri.

Asagidaki alanlar icin testler yazilmali:
- `helpers/file.go` - `ReadNTPConfFile`, `CheckFileExists`, `GetLargestFoldersConcurrent`
- `helpers/os.go` - `IsWindows`, `IsMacOS`
- `helpers/packages.go` - `GetHomeBrewVersion`, `GetPythonVersion`
- `types/file.go` - struct validation

Ornek test dosyasi yapisi:
```
helpers/
  file_test.go
  os_test.go
  packages_test.go
```

Ayrica CI/CD pipeline'ina test adimi eklenmeli.

---

## 7. CI/CD Iyilestirmeleri

### 7.1 `build-release.yml` - Copy-Paste Hatasi

Satir 65'te `darwin/amd64` build'i icin build tag yanlislikla `darwin-arm64` olarak ayarlanmis:

```yaml
release-mac-amd64:
    name: release darwin/amd64
    ...
        goarch: amd64
        build_tags: darwin-arm64  # HATALI: darwin-amd64 olmali
```

### 7.2 CI Pipeline Eksikligi

Yalnizca release pipeline'i mevcut. Asagidakiler eklenmeli:
- **Lint adimi**: `golangci-lint` ile kod kalitesi kontrolu
- **Test adimi**: `go test ./...`
- **Build kontrolu**: Her PR'da derleme testi
- **`go vet`**: Statik analiz

Ornek CI workflow:
```yaml
on: [push, pull_request]
jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: go vet ./...
      - run: go test ./...
      - uses: golangci/golangci-lint-action@v6
```

---

## 8. Windows Uyumluluk Kontrolu

Yalnizca `info` komutu Windows kontrolu yapiyor. `cpu` ve `largest-folders` komutlari Windows'ta sorunlu calisabilir. Kontrol merkezi bir yere alinmali:

```go
// root.go PersistentPreRun ile tum komutlarda kontrol
var rootCmd = &cobra.Command{
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        if helpers.IsWindows() {
            fmt.Println("This program is not compatible with Windows.")
            os.Exit(1)
        }
    },
}
```

Ayrica `os.Exit(0)` yerine `os.Exit(1)` kullanilmali - uyumsuzluk bir hata durumudur.

---

## 9. Versiyon Yonetimi

### 9.1 Redundant Sabitler

`constants/version.go`:
```go
const (
    VERSION = "0.0.9"
    MAJOR   = 0
    MINOR   = 0
    PATCH   = 9
)
```

`MAJOR`, `MINOR`, `PATCH` sabitleri hicbir yerde kullanilmiyor ve `VERSION` ile senkronizasyonu manuel.

### 9.2 Build-Time Version Injection

Versiyon bilgisi build sirasinda `ldflags` ile enjekte edilmeli:
```go
var Version = "dev" // ldflags ile override edilir
```

```bash
go build -ldflags "-X genius/constants.Version=1.0.0"
```

### 9.3 Go Isimlendirme Kurallari

`VERSION` yerine `Version` kullanilmali. Go'da ALL_CAPS convention'i yoktur; exported constant'lar icin PascalCase kullanilir.

---

## 10. Cikti Formatlama

### 10.1 Yapisal Cikti Destegi

Tum cikti `fmt.Println` ile duz metin olarak basilyor. JSON veya YAML formati secenegi eklenmeli:

```bash
genius info --output json
genius info --output yaml
genius info           # varsayilan: table/text
```

### 10.2 Renkli Cikti

`fatih/color` veya `charmbracelet/lipgloss` gibi kutuphaneler ile terminal ciktisi daha okunabilir hale getirilebilir.

### 10.3 `fmt.Println` Tutarsizligi

Bazi satirlarda `fmt.Println("Label: ", value)` seklinde fazla bosluk var (virgulden sonra Println otomatik bosluk ekler, string icindeki boslukla birlikte cift bosluk olusur).

---

## 11. Bagimliliklarin Guncellenmesi

### 11.1 gopsutil v3 -> v4

`github.com/shirou/gopsutil/v3` kullaniliyor, ancak v4 mevcut. v4'e gecis oneriliyor.

### 11.2 Gereksiz Agir Bagimlilk

`github.com/miekg/dns` kutuphanesi yalnizca `/etc/resolv.conf` dosyasini okumak icin kullaniliyor. Bu islem basit dosya okuma ve satir ayristirma ile de yapilabilir, bu agir DNS kutuphanesine gerek kalmadan.

---

## 12. Guvenlik

### 12.1 Harici Komut Calistirma

`helpers/packages.go` icinde `exec.Command` ile harici komutlar calistiriliyor. Her ne kadar mevcut durumda dogrudan bir risk olmasa da, ciktilarin sanitize edilmesi ve timeout eklenmesi oneriliyor:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
cmd := exec.CommandContext(ctx, "python3", "--version")
```

---

## 13. Dokumantasyon

### 13.1 README Yetersizligi

Mevcut README cok minimal. Eklenmesi gerekenler:
- Tum komutlarin dokumantasyonu (`info`, `cpu`, `largest-folders`)
- Ornek ciktilar (screenshots veya terminal ciktisi)
- Desteklenen platformlar listesi
- Katki rehberi (CONTRIBUTING.md)
- Badge'ler (build status, go version, license)

### 13.2 GoDoc Eksikligi

Bazi exported fonksiyonlarda GoDoc yorum satiri yok:
- `helpers/os.go`: `IsWindows()` ve `IsMacOS()` fonksiyonlarinda aciklama yok
- `types/file.go`: `NTPConfiguration` struct'inda alan aciklamalari yok

---

## 14. Performans

### 14.1 `GetLargestFolders` - O(n^2) Dosya Taramas

Non-concurrent `GetLargestFolders` fonksiyonunda ic ice `WalkDir` + `Walk` kullanimi, alt dizinlerin tekrar tekrar taranmasina neden olur. Her ne kadar concurrent versiyon tercih edilse de, bu fonksiyon ya duzeltilmeli ya da kaldirilmali.

### 14.2 `GetLargestFoldersConcurrent` - Goroutine Sinirlamasi

Ev dizininde cok fazla alt dizin varsa, sinirlandirma olmaksizin goroutine baslatiliyor. Bir `semaphore` veya `worker pool` patterni ile goroutine sayisi sinirlandirilmali:

```go
sem := make(chan struct{}, 10) // max 10 concurrent worker
for _, dir := range dirs {
    sem <- struct{}{}
    go func(d os.DirEntry) {
        defer func() { <-sem }()
        // ...
    }(dir)
}
```

---

## 15. Diger Kucuk Iyilestirmeler

| Dosya | Satir | Sorun | Oneri |
|---|---|---|---|
| `helpers/file.go:33` | `Errorf` | Error mesajinda nokta (`.`) var | Go convention'ina gore hata mesajlari kucuk harfle baslar ve noktasiz biter |
| `helpers/file.go:51` | `strings.Index` | `strings.Index` yerine `strings.Contains` kullanilmali | Daha okunabilir |
| `helpers/file.go:83` | return | Pointer to slice (`*[]Type`) donduruluyor | Go'da slice zaten referans tipidir, `[]Type` yeterli |
| `cmd/info.go:93-183` | output | Separator olarak `====` kullaniliyor | Sabit bir constant veya format fonksiyonu olusturulmali |
| `constants/version.go` | all | Dosya ismi | `constants` yerine Go idiomatic olarak `version` paketi veya root'ta basit bir degisken daha uygun |

---

## Oncelik Sirasi

| Oncelik | Alan | Etki |
|---|---|---|
| **Kritik** | Hata yonetimi (panic kaldirilmasi) | Uygulama kararliligi |
| **Kritik** | Type assertion guvenlik kontrolu | Cokme riski |
| **Kritik** | CI/CD build tag copy-paste hatasi | Yanlis build |
| **Yuksek** | Test eklenmesi | Kod guvenirligi |
| **Yuksek** | Ag arayuzu hardcode kaldirmasi | Tasinabilirlik |
| **Yuksek** | NTP parser duzeltmesi | Dogru calismasi |
| **Orta** | Kod organizasyonu (info.go parcalanmasi) | Bakimlanabilirlik |
| **Orta** | CI pipeline eklenmesi | Gelistirme sureci |
| **Orta** | JSON/YAML cikti destegi | Kullanilabilirlik |
| **Dusuk** | Go naming convention uyumu | Kod kalitesi |
| **Dusuk** | README iyilestirmesi | Dokumantasyon |
| **Dusuk** | Kullanilmayan kodun temizlenmesi | Kod temizligi |

---

*Bu rapor, projenin `v0.0.9` surumu uzerinden olusturulmustur.*
