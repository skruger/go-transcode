
CREATE TABLE IF NOT EXISTS source_asset (
    id integer not null primary key,
    filename varchar(256) not null unique,
    codec varchar(50) null,
    filesize int null,
    duration_time decimal(10,3) null,
    duration_frames int null,
    fps decimal(5,4) null,
    resolution_width int null,
    resolution_height int null
);

CREATE UNIQUE INDEX IF NOT EXISTS source_asset_filename ON source_asset (filename);

CREATE TABLE IF NOT EXISTS transcode_asset (
    id integer not null primary key,
    source int not null,
    filename varchar(256) not null,
    profile_name varchar(50),
    filesize int null,
    resolution_width int null,
    resolution_height int null,

    FOREIGN KEY(source) REFERENCES source_asset
);

CREATE TABLE IF NOT EXISTS asset_package (
    id integer not null primary key,
    source int not null,
    manifest_file varchar(256) not null,
    manifest_root varchar(256) not null,

    FOREIGN KEY(source) REFERENCES source_asset
);

CREATE TABLE IF NOT EXISTS asset_package_file (
    id integer not null primary key,
    package int not null,
    filename varchar(256) not null,

    FOREIGN KEY(package) REFERENCES asset_package
)
