# Go OAI

This originally started as an OAI scraper in go, but has turned into a quite specific scraper for a UKL web property. It can certainly be adapted, but its current form is quite reliant on the structure of that property.

# What's Missing?

Included is the all-ids.dip.txt which has a list of all id's in web property. You will need to create a run-dips.txt file for the dips you'd like to run. I also have it logging a collections-completed-log.txt so that you can see which id's have been ran. This accomplishes 2 things: The collection complete and run-dips, when diffed, should add up to the total within all-ids.dip.txt, making sure you've not missed something.

Also missing is a better way to store the files. Currently, all files are written to a single directory, which is quite annoying as you'll probably get hundreds of thousands of txt files. In the current form, I'm using a shell script that is something along the lines of this:

```bash
#!/bin/bash

# Check if batch name is provided
if [ -z "$1" ]; then
  echo "Usage: ./move.sh batch_name"
  exit 1
fi

BATCH_NAME="$1"
TARGET_DIR=".../go-oai/resources/${BATCH_NAME}"

# Create the target directory if it doesn't exist
mkdir -p "${TARGET_DIR}"

# Move the first 30000 files to the target directory
find .../go-oai/resources -maxdepth 1 -type f | head -n 30000 | xargs -I {} mv {} "${TARGET_DIR}"

echo "Moved 30000 files to ${TARGET_DIR}"
```

# What's good to know?

Right now, this all runs synchronously. Most of the time, the code is doing some analysis on the xml files, and then makes requests for new xml as necessary. This means that the amount of requests is relatively low per second.
Make your run-dips.txt or it won't run
Make sure you are regularly storing the txt and xml files in their own folders (I currently store them in 30000 batches)
