// package xulid adalah package internal untuk menghandle kebutuhan ulid aplikasi.
// Karena "github.com/oklog/ulid/v2" dengan field ulid.ULID adalah bertipe byte dan tidak memenuhi kriteria pada Value Scanner Database
// maka dibuatlah tipe data custom yang dapat diakses pada package ini
// xulid.ULID
// ulid pada aplikasi ini mengandalkan extensi postgresql untuk ulid : https://github.com/pksunkara/pgx_ulid
// ekstensi tersebut memungkinkan ulid disimpan di database menggunakan type data bytea seperti halnya uuid.UUID
// author : muchlis
package xulid
