# kakaPika 
Control TP-Link and Wemo "smart" devices.

<code>kakaPika</code> is a stand-alone binary that can control/query these devices.
# License 
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.txt)

# About the name: Cuckoos and Crows 

Cuckoos are an interesting species: [They](https://www.calacademy.org/explore-science/cuckoo-for-crows)  are [brood parasites](https://en.wikipedia.org/wiki/Brood_parasite).

A number of popular algorithms are named after them: 
* [Cuckoo Hashing](https://en.wikipedia.org/wiki/Cuckoo_hashing) - where the cuckoo chick pushes the other eggs or young out of the nest when it hatches. 
* [Cuckoo Search](https://en.wikipedia.org/wiki/Cuckoo_search) - if a host bird discovers the eggs are not their own, they will either throw these alien eggs away or simply abandon its nest and build a new nest elsewhere.


Here's a verse in sanskrit, [oldest living language](https://www.wisdomlib.org/sanskrit/quote/mss/subhashita-9283) on this interesting bird:

<pre>
काकः कृष्णः पिकः कृष्णः को भेदः पिककाकयोः । 
वसन्तसमये प्राप्ते काकः काकः पिकः पिकः ॥

kākaḥ kṛṣṇaḥ pikaḥ kṛṣṇaḥ ko bhedaḥ pikakākayoḥ |
vasantasamaye prāpte kākaḥ kākaḥ pikaḥ pikaḥ ||

കാകഃ കൃഷ്ണഃ, പികഃ കൃഷ്ണഃ,
കോ ഭേദഃ പികകാകയോഃ?
വസന്തകാലേ സമ്പ്രാപ്തേ
കാകഃ കാകഃ, പികഃ പികഃ! 
</pre>


In the wild, Cuckoos lay their eggs in Crow's nests. When chicks hatch, they all look and act the same. 

How can one tell them a part?

Well, in the Spring, crow begin to caw and the cuckoos of the brood, cuckoo!

## Crows, Cuckoos and the "smart" devices I have:

As the market is flooded with types and brands of "smarts," I found myself having to operate TP-Link's Kasa and Belkin's Wemo devices. They both come with their own "SmartPhone" app to control these devices. This was not acceptable and I wanted a single point of control for these devices. 

I essentially don't care about their brand. I should be able to control them as I please.

I wrote kakaPika to that effect. I am hoping this would be useful to others who have both these devices. 

# Installing
From the [Releases](https://github.com/evuraan/kakaPika/releases) page, download 
the pre-built binary, suitable to your operating system, to a folder of your choice. 
## Compiling
Optionally, if you want to compile your own binary, you would need to install [Go Programming Language](https://golang.org/).

Clone this repository and build as
<pre>
$ go build </pre>

# Usage 
How to operate and get status of your device:
## help
Use the <code>-h</code> for help:
<pre>
$ ./kakaPika -h
Usage: 
  -h  --help           print this usage and exit
  -d  --device         smartPlug hostname or address
  -c  --cmd            cmd to run: [on,off,stat]
  -v  --version        print version information and exit
</pre>

## status 
Use <code>stat</code> argument to find the current status of your device:
<pre>
$ ./kakaPika -d wemo_laundry -c stat 
device: wemo_laundry, cmd: stat
stat: off
.. later..
$ ./kakaPika -d wemo_laundry -c stat 
device: wemo_laundry, cmd: stat
stat: on
</pre>
## Control
Use <code>on</code> and <code>off</code> directives to control your device:
<pre>
$ ./kakaPika -d wemo_laundry -c off
device: wemo_laundry, cmd: off
off: off
</pre>
and,<pre>
$ ./kakaPika -d wemo_laundry -c on
device: wemo_laundry, cmd: on
on: on
</pre>
# Contribution
You are very welcome to contribute. If you have a new feature in mind, open an issue on github describing it. 

# References
* https://malayalam.usvishakh.net/blog/archives/191
