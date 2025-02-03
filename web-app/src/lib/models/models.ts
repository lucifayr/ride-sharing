import { User } from "./user";

export type RideGroup = {
  groupId: string;
  groupName: string;
  mainRide: Ride;
  mainSchedule: Schedule;
  members: User[];
};

export type ScheduledRideGroup = {
  secondaryRide: Ride;
  secondarySchedule: Schedule;
} & RideGroup;

type Ride = {
  rideId: string;
  destination: string;
  departurePoint: string;
};
export type Group = {
  groupId: string;
  name: string;
  description?: string;
  createdBy: string; // user id
  members: GroupMemeber[];
};

export type GroupMemeber = {
  userId: string;
  email: string;
  joinStatus: string;
};

export type Schedule = {
  monday: string;
  tuesday: string;
  wednesday: string;
  thrusday: string;
  friday: string;
  saturday: string;
  sunday: string;
};
