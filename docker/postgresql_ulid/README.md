compose postgresql_ulid digunakan untuk mendukung type ulid di postgresql agar dapat menyimpan string ulid sebagai byte dengan type spesifik 
```
"id" ulid DEFAULT gen_ulid() PRIMARY KEY
```

namun ternyata untuk memodifikasi postgresql baremetal menginstall `pgx_ulid` terkadang memiliki kendala. sehingga diputuskan untuk menggunakan ulid varchar VARCHAR(26).
