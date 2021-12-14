#!/bin/bash
jpg_files=`find ./www/content/post -type f -name "*.jpeg"`
for image in $jpg_files; do
  webp_image=`echo $image | sed -e "s/jpeg\$/webp/"`
  echo "$webp_image"
  magick convert -auto-orient -quality 85 -resize 60% $image $webp_image
  find . -name "*.md" | xargs sed -i "s/$(basename $image)/$(basename $webp_image)/g"
  rm $image
done
