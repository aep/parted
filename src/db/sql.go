package db

const SQLReadInboud = `
SELECT id,
       manufacturer,
       part_number,
       description,
       image,
       stock,
       used,
       order_number,
       barcode_id,
       "values",
       units,
       labels
FROM inventory
LEFT JOIN
(
	SELECT
    id_element,
	group_concat(specifications.label, ';;') as "labels",
	group_concat( specifications.value, ';;') as "values",
	group_concat( specifications.unit, ';;') as "units"
FROM specifications
GROUP BY id_element
) as s on s.id_element = inventory.id
 WHERE inventory.order_number = ?;
`

const SQLSchema = `
CREATE TABLE inventory (
    id INTEGER NOT NULL PRIMARY KEY,
    manufacturer VARCHAR NOT NULL DEFAULT '',
    part_number VARCHAR NOT NULL DEFAULT '',
    description VARCHAR NOT NULL DEFAULT '',
    image VARCHAR NOT NULL DEFAULT '',
    stock INTEGER NOT NULL DEFAULT 0,
    used INTEGER NOT NULL DEFAULT 0,
    order_number VARCHAR NOT NULL DEFAULT '',
    barcode_id INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE specifications (
    id_element INTEGER NOT NULL,
    label VARCHAR NOT NULL DEFAULT '',
    value VARCHAR NOT NULL DEFAULT '',
    unit VARCHAR NOT NULL DEFAULT '',
    FOREIGN KEY(id_element) REFERENCES inventory(id)
);
`

const SQLInsertItem = `
INSERT INTO inventory (
	manufacturer,
	part_number,
	description,
	image,
	stock,
	order_number,
	barcode_id
)
VALUES (?, ?, ?, ?, ?, ?, ?);
`

const SQLInsertAttr = `
INSERT INTO specifications (id_element, label, value, unit)
VALUES (?, ?, ?, ?);
`

const SQLReadItems = `
SELECT id,
       manufacturer,
       part_number,
       description,
       image,
       stock,
       used,
       order_number,
       barcode_id,
       "values",
       units,
       labels
FROM inventory
LEFT JOIN
(
	SELECT
    id_element,
	group_concat(specifications.label, ';;') as "labels",
	group_concat( specifications.value, ';;') as "values",
	group_concat( specifications.unit, ';;') as "units"
FROM specifications
GROUP BY id_element
) as s on s.id_element = inventory.id;
`

const DeleteInbound = `
DELETE FROM inventory WHERE inventory.order_number = ?
RETURNING id;
`

const SQLDeleteAttr = `
DELETE FROM specifications WHERE id_element = ?;
`

const SQLidsFromInbound = `
SELECT id
FROM inventory WHERE order_number = ?;
`
