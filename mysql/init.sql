SET SQL_MODE = "TRADITIONAL";
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET SQL_MODE = "ONLY_FULL_GROUP_BY";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";
SET NAMES "utf8mb4" COLLATE "utf8mb4_bin";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `client+twitterclone+naro`
--

CREATE DATABASE `twitterclone`;
USE `twitterclone`;
-- --------------------------------------------------------

--
-- テーブルの構造 `tweets`
--

CREATE TABLE `tweets` (
  `id` CHAR(36) NOT NULL PRIMARY KEY,
  `user_id` VARCHAR(32) NOT NULL,
  `tweet_body` VARCHAR(256) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- --------------------------------------------------------

--
-- テーブルの構造 `users`
--

CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `hashed_pass` varchar(256) NOT NULL,
  PRIMARY KEY (`id`)
);

-- --------------------------------------------------------

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
