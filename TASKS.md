This is a copy paste skeleton we use for tools in the company. Nothing is fixed, feel free to change anything.

The project we're building here is an internal tool for stock keeping of electronics parts. It has these main building blocks:

1. warehouse inventory

a list of items in the warehouse and their specs.

 - manufacturer
 - manufacturer part number
 - description
 - image (optional)
 - category, such as resistor, capacitor, mcu, .._
 - spec , they're dynamic key/value fields, such as capacitance=6uF, height=1mm
 - barcode with its database id and the manufacturer part number
 - the amount available in inventory
 - the amount used by projects
 - the inbound order nr this was purchased from


2. warehouse inbound

an inbound order is a package from a distributor with bags in them that contain components.
This page creates a new inbound order with an order reference, and multiple item that are in the package.
The bags will be scanned with a barcode scanner that emulates a keyboard, pressing the return key on every new item. 

The website must catch the enter key, search for what was entered on element14, and replace the input line with a multi-selection for the matching components.
then a new row must appear to scan the next item.

Then the order is completed, all rows are created as inventory items.
If the item was already in the inventory database, we create a new one ANYWAY, because these items are not mixable goods.


3. bill of material assignments

a bill of materials is uploaded as csv and then matched against the inventory using text search on ALL the fields, including the dynamic spec
for every row, a multiselect box will appear will all the matched components, and the engineer has to select the correct one.
once finished, for each item in the bom, the "used" counter in the inventory is increased.

later, we will want to edit, and delete BOMs, changing the counter of the inventory again.




