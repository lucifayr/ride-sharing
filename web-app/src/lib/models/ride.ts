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
  status: string;
  schedule: RideSchedule | null;
  participants: RideParticipant[];
};

export type RideSchedule = {
  unit: string;
  interval: number;
  weekdays: string[] | null;
};

export type RideParticipant = {
  userId: string;
  email: string;
};
