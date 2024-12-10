import { useForm } from "@tanstack/react-form";
import { ScheduledRideGroup } from "../models/models";

type ScheduledRideFormType = {
  groupName: string;
  destination: string;
  departurePoint: string;
};

type ScheduleFormType = {
  monday: Date;
  tuesday: Date;
  wednesday: Date;
  thrusday: Date;
  friday: Date;
  saturday: Date;
  sunday: Date;
};

export default function ScheduledRideForm() {
  const form = useForm<ScheduledRideFormType>({
    defaultValues: {
      groupName: "",
      destination: "",
      departurePoint: "",
    },
    onSubmit: async ({ value }) => {
      console.log(value);
    },
  });
  return (
    <>
      <form.Field
        name="groupName"
        children={(field) => {
          return (
            <div>
              <label className="font-bold">Group Name:</label>
              <input
                placeholder="Gruppenname"
                id={field.name}
                name={field.name}
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => {
                  field.handleChange(e.target.value);
                }}
              />
            </div>
          );
        }}
      ></form.Field>
    </>
  );
}
