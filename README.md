<!DOCTYPE html>
<html lang="en">ll
<body>

<h1>Country Info API</h1>
<p>
  Welcome to the Country Info API! This service provides endpoints to retrieve
  a country’s core details, its population, and a quick status overview of the API itself.
</p>

<hr />

<h2>Base URL</h2>
<p>
  <strong>Production:</strong> 
  <code>https://prog20052025-bsio.onrender.com</code>
</p>
<p>
  All endpoints below are relative to the base URL.
</p>

<hr />

<h2>Endpoints</h2>

<h3>1. <code>GET /countryinfo/v1/info/<em>{countryCode}</em></code></h3>
<p>
  Retrieves a JSON object of detailed country information:
</p>
<ul>
  <li>Country name</li>
  <li>Continents</li>
  <li>Population</li>
  <li>Languages</li>
  <li>Borders</li>
  <li>Capital</li>
  <li>Flag</li>
  <li>List of major cities (limited by <code>limit</code> parameter)</li>
</ul>

<p><strong>Path parameter:</strong></p>
<ul>
  <li><code>countryCode</code> – A 2-letter ISO code (e.g. <code>US</code>, <code>NO</code>, <code>FR</code>).</li>
</ul>

<p><strong>Query parameter (optional):</strong></p>
<ul>
  <li>
    <code>limit</code> – Integer for how many cities to return 
    (<em>default</em> = 10).
  </li>
</ul>

<p><strong>Example Request:</strong></p>
<pre><code>GET /countryinfo/v1/info/US?limit=5
Host: prog20052025-bsio.onrender.com
</code></pre>

<p><strong>Example Response:</strong></p>
<pre><code>{
  "Name": "United States",
  "Continents": ["North America"],
  "Population": 331002651,
  "Languages": ["English"],
  "Borders": ["CAN", "MEX"],
  "Capital": "Washington D.C.",
  "Cities": ["New York", "Los Angeles", "Chicago", "Houston", "Phoenix"],
  "Flag": "https://example.com/flags/US.png"
}
</code></pre>

<hr />


<h3>2. <code>GET /countryinfo/v1/population/<em>{countryCode}</em></code></h3>
<p>
  Returns a JSON object containing only the population information for the specified country.
</p>

<p><strong>Path parameter:</strong></p>
<ul>
  <li><code>countryCode</code> – A 2-letter ISO code (e.g. <code>US</code>, <code>NO</code>, <code>FR</code>).</li>
</ul>

<p><strong>Example Request:</strong></p>
<pre><code>GET /countryinfo/v1/population/US
Host: prog20052025-bsio.onrender.com
</code></pre>

<p><strong>Example Response:</strong></p>
<pre><code>{
  "Country": "United States",
  "Population": 331002651
}
</code></pre>

<hr />


<h3>3. <code>GET /countryinfo/v1/status</code></h3>
<p>
  Returns a JSON object with general information about the API's status (e.g., version, uptime, or a simple "OK" message).
</p>

<p><strong>Example Request:</strong></p>
<pre><code>GET /countryinfo/v1/status
Host: prog20052025-bsio.onrender.com
</code></pre>

<p><strong>Example Response:</strong></p>
<pre><code>{
  "Status": "OK",
  "Version": "1.0.0",
  "Uptime": "72h"
}
</code></pre>

<hr />

<h2>Error Handling</h2>
<p>
  The API returns errors in a JSON structure with <code>StatusCode</code> and <code>Message</code> fields, for example:
</p>
<pre><code>{
  "Message": "Invalid country code. Country code should only consist of 2 characters.",
  "StatusCode": 400
}
</code></pre>

<p><strong>Common error codes:</strong></p>
<ul>
  <li>
    <strong>400 Bad Request:</strong> The country code is invalid
    (e.g. not exactly 2 letters), or the external data source indicates
    a missing/invalid country.
  </li>
  <li>
    <strong>404 Not Found:</strong> The endpoint does not exist
    (e.g. a typo in the URL).
  </li>
  <li>
    <strong>405 Method Not Allowed:</strong> Attempted an HTTP method other than GET on these endpoints.
  </li>
  <li>
    <strong>500 Internal Server Error:</strong> An unexpected error occurred
    on the server or with upstream data sources.
  </li>
</ul>


</body>
</html>
