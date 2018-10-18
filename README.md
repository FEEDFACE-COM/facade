tcpdump -l -i pppoe0 -a  -t "tcp[tcpflags] & (tcp-syn) != 0 and tcp[tcpflags]&(tcp-ack) == 0" | fcd
tcpdump -l  -i vlan5 -t   dst port 53  | sed -u 's/.hq.feedface.com//g;s/.tsaoyidi.com//g'  | fcd