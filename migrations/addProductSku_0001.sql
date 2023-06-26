IF NOT EXISTS( SELECT NULL FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = 'Products'AND table_schema = 'store' AND column_name = 'sku') THEN ALTER TABLE `Products` ADD `sku` VARCHAR(255);
END IF;