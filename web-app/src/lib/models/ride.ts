export type RideEvent = {
  rideId: string;
  rideEventId: string;
  locationFrom: string;
  locationTo: string;
  tackingPlaceAt: string;
  createdBy: string;
  createdByEmail: string;
  driverId: string;
  driverEmail: string;
  transportLimit: number;
  schedule?: RideSchedule;
};

export type RideSchedule = {
  unit: string;
  interval: number;
  weekdays?: string[];
};
