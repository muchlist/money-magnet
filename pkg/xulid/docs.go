// package xulid adalah paket internal untuk menghandle kebutuhan ulid pada keseluruhan aplikasi ini.
// sebab, "github.com/oklog/ulid/v2" dengan field ulid.ULID berupa byte tidak memenuhi kriteria pada Value Scanner.
// xulid.ULID mengimplementasi Value, Scanner, Marshalling yang sudah disesuaikan.
// database aplikasi ini mengandalkan extensi postgresql untuk ulid : https://github.com/pksunkara/pgx_ulid
// ekstensi tersebut memungkinkan ulid disimpan di database menggunakan type data bytea seperti halnya uuid.UUID
package xulid
