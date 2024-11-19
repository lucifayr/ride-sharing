import { RideGroup, ScheduledRideGroup } from "./models";
import { User } from "./user";

export type UniversalSocketMessage = {
  type: string;
} & WSScheduledRideGroupUpdate &
  WSRideGroupUpdate;

export type ChatMessage = {
  timestamp: Date;
  content: string;
  user: User;
  group: RideGroup;
};

type WSScheduledRideGroupUpdate = {
  group: ScheduledRideGroup;
};

type WSRideGroupUpdate = {
  group: RideGroup;
};
