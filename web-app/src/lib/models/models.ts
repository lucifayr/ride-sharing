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

export type GroupMessage = {
  groupId: string;
  messageId: string;
  content: string;
  sentBy: string; // user id
  sentByEmail: string;
  createdAt: string;
  repliesTo?: string;
};
