create database if not exists WirelessPA;

use WirelessPA;

create table if not exists LO (
  recordType char(2),
  uniqueSystemIdentifier numeric(9, 0),
  ulsFileNumber char(14),
  ebfNumber varchar(30),
  callSign char(10),
  locationActionPerformed char(1),
  locationTypeCode char(1),
  locationClassCode char(1),
  locationNumber integer,
  siteStatus char(1),
  correspondingFixedLocation integer,
  locationAddress varchar(80),
  locationCity char(20),
  locationCountyBoroughParish varchar(60),
  locationState char(2),
  radiusOfOperation numeric(5, 1),
  areaOfOperationCode char(1),
  clearanceIndicator char(1),
  groundElevation numeric(7, 1),
  latitudeDegrees integer,
  latitudeMinutes integer,
  latitudeSeconds numeric(3, 1),
  latitudeDirection char(1),
  longitudeDegrees integer,
  longitudeMinutes integer,
  longitudeSeconds numeric(3, 1),
  longitudeDirection char(1),
  maxLatitudeDegrees integer,
  maxLatitudeMinutes integer,
  maxLatitudeSeconds numeric(3, 1),
  maxLatitudeDirection char(1),
  maxLongitudeDegrees integer,
  maxLongitudeMinutes integer,
  maxLongitudeSeconds numeric(3, 1),
  maxLongitudeDirection char(1),
  nepa char(1),
  quietZoneNotificationDate date,
  towerRegistrationNumber char(10),
  heightOfSupportStructure numeric(7, 1),
  overallHeightOfStructure numeric(7, 1),
  structureType char(7),
  airportId char(4),
  locationName char(20),
  unitsHandHeld integer,
  unitsMobile integer,
  unitsTempFixed integer,
  unitsAircraft integer,
  unitsItinerant integer,
  statusCode char(1),
  statusDate date,
  earthStationAgreement char(1),
  primary key (uniqueSystemIdentifier)
);

create table if not exists FR (
  recordType char(2),
  uniqueSystemIdentifier numeric(9, 0),
  ulsFileNumber char(14),
  ebfNumber varchar(30),
  callSign char(10),
  frequencyActionPerformed char(1),
  locationNumber integer,
  antennaNumber integer,
  classStationCode char(4),
  opAltitudeCode char(2),
  frequencyAssigned numeric(16, 8),
  frequencyUpperBand numeric(16, 8),
  frequencyCarrier numeric(16, 8),
  timeBeginOperations integer,
  timeEndOperations integer,
  powerOutput numeric(15, 3),
  powerErp numeric(15, 3),
  tolerance numeric(6, 5),
  frequencyIndicator char(1),
  status char(1),
  eirp numeric(7, 1),
  transmitterMake varchar(25),
  transmitterModel varchar(25),
  autoTransmitterPowerControl char(1),
  numberOfUnits integer,
  numberOfPagingReceivers integer,
  frequencyNumber integer,
  statusCode char(1),
  statusDate date,
  primary key (uniqueSystemIdentifier, antennaNumber)
);

-- F2 is currently unused
create table if not exists F2 (
  recordType char(2),
  uniqueSystemIdentifier numeric(9, 0),
  ulsFileNumber char(14),
  ebfNumber varchar(30),
  callSign char(10),
  additionalFrequencyInfoActionPerformed char(1),
  locationNumber integer,
  antennaNumber integer,
  frequencyNumber integer,
  frequencyAssigned numeric(16, 8),
  frequencyUpperBand numeric(16, 8),
  offset char(3),
  frequencyChannelBlock char(4),
  equipmentClass char(2),
  minimumPowerOutput numeric(15, 3),
  dateFirstUsed date,
  statusCode char(1),
  statusDate date,
  protocolRestrictedOrUnrestricted char(1),
  primary key (uniqueSystemIdentifier, antennaNumber)
);