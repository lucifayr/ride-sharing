import { User } from "./user";

type RideGroup = {
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

export type Ride = {
  rideId: string;
  destination: string;
  departurePoint: string;
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
