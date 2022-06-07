CREATE TABLE `users` (
                         `id` int(11) NOT NULL AUTO_INCREMENT,
                         `role` tinyint(1) DEFAULT NULL COMMENT '0 = Admin, 1 = Coustomer, 2 = Pharmacist',
                         `name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                         `email` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                         `gender` tinyint(1) DEFAULT '0',
                         `password` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
                         `picture` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
                         `country_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
                         `phone` bigint(20) DEFAULT NULL,
                         `stripe_id` text COLLATE utf8_unicode_ci,
                         `is_test` tinyint(1) DEFAULT '1',
                         `is_verify` tinyint(1) DEFAULT '1',
                         `is_active` tinyint(1) DEFAULT '1',
                         `is_delete` tinyint(1) DEFAULT '0',
                         `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                         `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci

CREATE TABLE `stores` (
                          `id` int(11) NOT NULL AUTO_INCREMENT,
                          `store_name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                          `store_image` text COLLATE utf8_unicode_ci,
                          `license_id` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                          `pharmacy_id` char(36) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
                          `is_test` tinyint(1) DEFAULT '1',
                          `is_verify` tinyint(1) DEFAULT '1',
                          `is_active` tinyint(1) DEFAULT '1',
                          `is_delete` tinyint(1) DEFAULT '0',
                          `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                          `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          `user_id` int(11) DEFAULT NULL,
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `license_id` (`license_id`),
                          UNIQUE KEY `pharmacy_id` (`pharmacy_id`),
                          KEY `userId` (`user_id`),
                          CONSTRAINT `stores_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci

CREATE TABLE `addresses` (
                             `id` int(11) NOT NULL AUTO_INCREMENT,
                             `primary_address` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
                             `addition_address_info` text COLLATE utf8_unicode_ci NOT NULL,
                             `address_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 = Home, 1  = Office, 2 = Location',
                             `latitude` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '21.228125',
                             `longitude` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '72.833771',
                             `is_select` tinyint(1) NOT NULL DEFAULT '0',
                             `is_test` tinyint(1) DEFAULT '1',
                             `is_delete` tinyint(1) DEFAULT '0',
                             `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                             `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             `user_id` int(11) DEFAULT NULL,
                             PRIMARY KEY (`id`),
                             KEY `userId` (`user_id`),
                             CONSTRAINT `addresses_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci

CREATE TABLE `prescriptions` (
                                 `id` int(11) NOT NULL AUTO_INCREMENT,
                                 `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
                                 `text_note` text COLLATE utf8_unicode_ci NOT NULL,
                                 `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 = pending, 1  = completed, 2 = Rejected',
                                 `is_test` tinyint(1) DEFAULT '1',
                                 `is_verify` tinyint(1) DEFAULT '1',
                                 `is_active` tinyint(1) DEFAULT '1',
                                 `is_delete` tinyint(1) DEFAULT '0',
                                 `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                 `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                 `user_id` int(11) DEFAULT NULL,
                                 PRIMARY KEY (`id`),
                                 KEY `userId` (`user_id`),
                                 CONSTRAINT `prescriptions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=92 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci

CREATE TABLE `prescription_images` (
                                       `id` int(11) NOT NULL AUTO_INCREMENT,
                                       `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
                                       `url` text COLLATE utf8_unicode_ci NOT NULL,
                                       `type` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
                                       `is_test` tinyint(1) DEFAULT '1',
                                       `is_verify` tinyint(1) DEFAULT '1',
                                       `is_active` tinyint(1) DEFAULT '1',
                                       `is_delete` tinyint(1) DEFAULT '0',
                                       `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                       `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                       `prescription_id` int(11) DEFAULT NULL,
                                       PRIMARY KEY (`id`),
                                       KEY `prescriptionId` (`prescription_id`),
                                       CONSTRAINT `prescription_images_ibfk_1` FOREIGN KEY (`prescription_id`) REFERENCES `prescriptions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=90 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci

CREATE TABLE `medicines` (
                             `id` int(11) NOT NULL AUTO_INCREMENT,
                             `name` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
                             `is_test` tinyint(1) DEFAULT '1',
                             `is_verify` tinyint(1) DEFAULT '1',
                             `is_active` tinyint(1) DEFAULT '1',
                             `is_delete` tinyint(1) DEFAULT '0',
                             `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                             `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             `prescription_id` int(11) DEFAULT NULL,
                             PRIMARY KEY (`id`),
                             KEY `prescriptionId` (`prescription_id`),
                             CONSTRAINT `medicines_ibfk_1` FOREIGN KEY (`prescription_id`) REFERENCES `prescriptions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=108 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci

CREATE TABLE `quotes` (
                          `id` int(11) NOT NULL AUTO_INCREMENT,
                          `store_name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                          `price` decimal(5,2) NOT NULL,
                          `text_note` text COLLATE utf8_unicode_ci NOT NULL,
                          `is_test` tinyint(1) DEFAULT '1',
                          `is_verify` tinyint(1) DEFAULT '1',
                          `is_active` tinyint(1) DEFAULT '1',
                          `is_delete` tinyint(1) DEFAULT '0',
                          `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                          `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          `store_id` int(11) DEFAULT NULL,
                          `prescription_id` int(11) DEFAULT NULL,
                          PRIMARY KEY (`id`),
                          KEY `storeId` (`store_id`),
                          KEY `prescriptionId` (`prescription_id`),
                          CONSTRAINT `quotes_ibfk_1` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
                          CONSTRAINT `quotes_ibfk_2` FOREIGN KEY (`prescription_id`) REFERENCES `prescriptions` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci