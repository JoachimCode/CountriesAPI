<!DOCTYPE html>
<html lang="en">
<head>

<h1>Country Info API</h1>
<p>
  Welcome to the Country Info API! This service provides convenient endpoints to retrieve
  a country’s core information, its flag, and a list of major cities.
</p>

<hr>

<h2>Base URL</h2>
<p>
  <strong>Local (example):</strong> <code>http://localhost:8080</code><br/>
  <strong>Production (example):</strong> <code>https://api.example.com</code>
</p>
<p>
  All endpoints below are relative to the base URL.
</p>

<hr>

<h2>Endpoints</h2>

<h3>1. <code>GET /info/<em>{countryCode}</em></code></h3>
<p>
  Retrieves a JSON object of country information, including:
</p>
<ul>
  <li>Country name</li>
  <li>Continents</li>
  <li>Population</li>
  <li>Languages</li>
  <li>Borders</li>
  <li>Capital</li>
  <li>Flag</li>
  <li>List of cities</li>
</ul>

<p><strong>Path parameter:</strong></p>
<ul>
  <li><code>countryCode</code> – An ISO 2-letter country code (e.g., <code>US</code>, <code>NO</code>, <code>FR</code>).</li>
</ul>

<p><strong>Query parameter (optional):</strong></p>
<ul>
  <li><code>limit</code> – Integer indicating how many cities to return. Defaults to 10.</li>
</ul>

<p><strong>Example Request:</strong></p>
<pre><code>GET /info/US?limit=5
Host: api.example.com
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

<hr>

<h2>Response Fields</h2>
<ul>
  <li><strong>Name</strong>: Common name of the country.</li>
  <li><strong>Continents</strong>: Array of continents the country spans.</li>
  <li><strong>Population</strong>: Total population.</li>
  <li><strong>Languages</strong>: Array of official/common languages.</li>
  <li><strong>Borders</strong>: Array of country codes for neighboring countries.</li>
  <li><strong>Capital</strong>: Name of the capital city.</li>
  <li><strong>Cities</strong>: List of major cities in that country (limited by <code>limit</code>).</li>
  <li><strong>Flag</strong>: URL or base64 data representing the country’s flag.</li>
</ul>

<hr>

<h2>Error Handling</h2>
<p>The API returns errors in a JSON structure with <code>StatusCode</code> and <code>Message</code> fields. Common scenarios include:</p>

<ul>
  <li>
    <strong>400 Bad Request:</strong> Invalid country code (e.g., fewer/more than 2 letters),
    or the external data source indicates a missing/invalid country.
  </li>
  <li>
    <strong>404 Not Found:</strong> You may see this if the specified endpoint does not exist.
  </li>
  <li>
    <strong>405 Method Not Allowed:</strong> You used an unsupported HTTP method
  </li>
  <li>
    <strong>500 Internal Server Error:</strong> An unexpected issue occurred on the server side
    or with upstream APIs.
  </li>
</ul>

<p><strong>Example Error Response:</strong></p>
<pre><code>{
  "Message": "Invalid country code. Country code should only consist of 2 characters.",
  "StatusCode": 400
}
</code></pre>

</body>
</html>
