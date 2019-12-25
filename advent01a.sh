cat advent01a.txt | sed -e "s/\(.*\)/(\1 \/ 3 - 2)/" | paste -s -d + - | bc
