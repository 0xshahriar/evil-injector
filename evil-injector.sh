#!/bin/bash

# ANSI color codes
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to check if "evil.com" is present in the response headers
check_for_evil() {
    while IFS= read -r line; do
        if [[ "$line" =~ "evil.com" ]]; then
            echo -e "${RED}Host header injection has been found on $domain${NC}"
            return 0
        fi
    done
    return 1
}

echo ""
echo "    ______      _ __   ____        _           __"
echo "   / ____/   __(_) /  /  _/___    (_)__  _____/ /_____  _____"
echo "  / __/ | | / / / /   / // __ \  / / _ \/ ___/ __/ __ \/ ___/"
echo " / /___ | |/ / / /  _/ // / / / / /  __/ /__/ /_/ /_/ / /"
echo "/_____/ |___/_/_/  /___/_/ /_/_/ /\___/\___/\__/\____/_/"
echo "                            /___/"
echo "By 0xShahriar"
echo ""

# Read the domain list file provided as an argument
if [ $# -ne 1 ]; then
    echo "Usage: $0 <domain_list_file>"
    exit 1
fi

domain_file="$1"

# Check if the domain list file exists
if [ ! -f "$domain_file" ]; then
    echo "Domain list file not found: $domain_file"
    exit 1
fi

# Get total number of domains
total_domains=$(wc -l < "$domain_file")
current_domain=0

# Read each domain from the text file
while IFS= read -r domain; do
    ((current_domain++))
    # Add "https://" before the domain
    url="https://$domain"
    
    # Perform silent curl with custom headers (X-Forwarded-Host) and check for "evil.com" in response headers
    echo -ne "Processing domain $current_domain/$total_domains\r"
    if curl -s -S -i -H "X-Forwarded-Host: evil.com" "$url" | check_for_evil; then
        continue
    fi
done < "$domain_file"

echo -e "\nProcessing completed."
