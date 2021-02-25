# FCC Wireless Public Access Database Contianer

This repo provides a mysql container that loads and serves the FCC's wireless public access datafiles in a relational database. The FCC's ftp directory is found here: ftp://wirelessftp.fcc.gov/pub/uls/complete/.

## Configuration

For now, we are only loading `l_cell.zip` and creating tables for it's `LO.dat` and `FR.dat`. The entire ftp directory is ~10 GB.

To load other .dat files into sql tables, add another json key/value to `src/config.json`. Then, add another `CREATE TABLE` entry for the .dat file that you want in `src/init.sql`.

You can extract the names and types for new tables from `src/uls_data_file_formats.pdf` by isolating your pdf page using https://www.ilovepdf.com/ and OCRing the pdf table to excel format with https://www.onlineocr.net/

## Requirements

You will need [Docker installed](https://docs.docker.com/get-docker/).

You do not need Go installed - this uses the compiled binary.

## Usage

To host this data in mysql, pull down a mysql docker image:

```bash
make dockerRunDB
```

To download data from the FCC and seed into mysql running in docker:

```bash
make seed
```

You should be able to connect to the database running in the container with the following:

```
host: localhost
username: root
password: root
database: WirelessPA
port: 3306
```

## Resources

Documentation on the FCC files can be found in the following links

https://www.fcc.gov/sites/default/files/pubacc_uls_code_def_02162017.txt

https://www.fcc.gov/sites/default/files/pubacc_tbl_abbr_names_08212007.pdf

## Compatibility

This should work on all GNU/Unix systems.

If you run Windows you can probably figure out how to run the commands in the `Makefile`. You will need bash.
