#!/bin/bash
jpg_files=`find ./www/content/post -type f -name "*.jpg"`
for image in $jpg_files; do
  webp_image=`echo $image | sed -e "s/jpg\$/webp/"`
  echo "$webp_image"
  magick convert -quality 85 -resize 70% $image $webp_image
  find . -name "*.md" | xargs sed -i "s/$(basename $image)/$(basename $webp_image)/g"
  rm $image
done
