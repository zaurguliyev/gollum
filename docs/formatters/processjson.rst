ProcessJSON
===========

ProcessJSON is a formatter that allows modifications to fields of a given JSON message.
The message is modified and returned again as JSON.


Parameters
----------

**ProcessJSONDataFormatter**
  ProcessJSONDataFormatter formatter that will be applied before ProcessJSONDirectives are processed.

**ProcessJSONGeoIPFile**
  ProcessJSONGeoIPFile defines a GeoIP file to load.
  This enables the "geoip" directive.
  If no file is loaded IPs will not be resolved.
  Files can be found e.g. at http://dev.maxmind.com/geoip/geoip2/geolite2/.

**ProcessJSONDirectives**
  ProcessJSONDirectives defines the action to be applied to the json payload.
  Directives are processed in order of appearance.
  The directives have to be given in the form of key:operation:parameters, where operation can be one of the following.
   * `split:<string>{:<key>:<key>:...}` Split the value by a string and set the    resulting array elements to the given fields in order of appearance. 
   * `replace:<old>:<new>` replace a given string in the value with a new one  * `trim:<characters>` remove the given characters (not string!) from the start    and end of the value  * `rename:<old>:<new>` rename a given field  * `remove{:<string>:<string>...}` remove a given field. If additional parameters are    given, an array is expected. Strings given as additional parameters will be removed    from that array  * `pick:<key>:<index>:<name>` Pick a specific index from an array and store it    in a new field. 
   * `time:<read>:<write>` read a timestamp and transform it into another    format  * `unixtimestamp:<read>:<write>` read a unix timestamp and transform it into another    format. valid read formats are s, ms, and ns. 
   * `flatten{:<delimiter>}` create new fields from the values in field, with new    fields named field + delimiter + subfield. Delimiter defaults to ".". Removes the original field. 
   * `agent:<key>{:<field>:<field>:...}` Parse the value as a user agent string and    extract the given fields into <key>_<field>. ("ua:agent:browser:os" would create the new fields "ua_browser" and "ua_os"). Possible values are: "mozilla", "platform", "os", "localization", "engine",    "engine_version", "browser", "version". 
   * `ip` Parse the field as an array of strings and remove all values that cannot    be parsed as a valid IP. Single-string fields are supported, too, but will be    converted to an array. 
   * `geoip:{<field>:<field>:...}` like agent this directive will analyse an IP string    via geoip and produce new fields. Possible values are: "country", "city", "continent", "timezone", "proxy", "location". 

**ProcessJSONTrimValues**
  ProcessJSONTrimValues will trim whitspaces from all values if enabled.
  Enabled by default.

Example
-------

.. code-block:: yaml

	- "stream.Broadcast":
	    Formatter: "format.ProcessJSON"
	    ProcessJSONDataFormatter: "format.Forward"
	    ProcessJSONGeoIPFile: ""
	    ProcessJSONDirectives:
	        - "host:split: :host:@timestamp"
	        - "@timestamp:time:20060102150405:2006-01-02 15\\:04\\:05"
	        - "error:replace:°:\n"
	        - "text:trim: \t"
	        - "foo:rename:bar"
	        - "foobar:remove"
	        - "array:pick:0:firstOfArray"
	        - "array:remove:foobar"
	        - "user_agent:agent:browser:os:version"
	        - "client:geoip:country:city:timezone:location"
	    ProcessJSONTrimValues: true
