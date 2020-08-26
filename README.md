
# FACADE by FEEDFACE.COM


    alias fcd='nc localhost 4045' 



# RFCs

    facade recv lines -w=72 -h=20 -vert=crawl -smooth &
    curl -L https://tools.ietf.org/rfc/rfc791.txt | while read line; \do echo $line; sleep .5; done | fcd
    
    
# Manpages    
    
    export MANWIDTH=50 MANPAGER=cat
    facade recv lines -w 50 -h 20 -vert roll &
    man ssh | while read line; do echo $line; sleep 1; done | fcd



# Access Logs

    facade recv lines -w 120 -h 8 -vert disk -mask mask &
    tail -f /var/log/nginx/access.log | fcd
    


