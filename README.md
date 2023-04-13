This little application will display some stats of your growatt power inverter using the API at https://server.growatt.com.

# Installation

```bash
go install github.com/patrickhener/growatt-display@latest
```

# Convert your password

There is a build in mode to convert your password to the format used by the API.

```bash
growatt-display -mode genhash
Enter your password: **************************
Provide this hash with -password <hash>: e530a41246ddcabf8d7f5c1b2bcfa0d1
```

You can now use the hashed password to login

# Display once

Only one shot login and get stats.

```bash
growatt-display -username <your-username> -password <hashed-password>
Login successful

Plant 'plant-name':
	Total Energy Today: 4.3 kWh
	Total Energy All Time: 5.1 kWh
	Total Co² reduction: 2.04 kg

Data collectors:

Collector 'collector-name'
	Current Power: 0.53kW
```

# Display loop

This command will keep your display "open" and refresh it after timeout (milliseconds) is over.

```bash
growatt-display -username <your-username> -password <hashed-password> -loop -timeout 60000
```

```
Plant 'plant-name':
	Total Energy Today: 4.3 kWh
	Total Energy All Time: 5.1 kWh
	Total Co² reduction: 2.04 kg

Data collectors:

Collector 'collector-name'
	Current Power: 0.53kW
```