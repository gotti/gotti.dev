#!/bin/bash
jpg_files=`find ./www/content/post -type f -name "*.jpg"`
for image in $jpg_files; do
  webp_image=`echo $jpg_files | sed -e "s/jpg^/webp/"`
  magick $jpg_files $webp_image
  rm $image
done
