module Main where

import Data.Function ( (&) )
import Data.List ( intercalate )

hello :: String -> String
hello s =
  "Hello, " ++ s ++ "."

main :: IO ()
main =
  map hello [ "artichoke", "alcachofa" ] & intercalate "\n" & putStrLn

-- Alcachofa, if you were wondering, is artichoke in Spanish.
