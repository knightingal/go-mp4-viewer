CREATE TABLE `mp4_base_dir` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `dir_path` varchar(100) DEFAULT NULL,
  `url_prefix` varchar(200) DEFAULT NULL,
  `api_version` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `video_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `video_file_name` varchar(200) DEFAULT NULL,
  `cover_file_name` varchar(200) DEFAULT NULL,
  `dir_path` varchar(100) DEFAULT NULL,
  `base_index` int(11) DEFAULT NULL,
  `rate` int(4) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `miss_match_video_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `video_file_name` varchar(200) DEFAULT NULL,
  `dir_path` varchar(100) DEFAULT NULL,
  `base_index` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE video_info ADD designation_char varchar(100) NULL;
ALTER TABLE video_info ADD designation_num varchar(100) NULL;
