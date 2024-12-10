import { useForm } from "@tanstack/react-form";

export default function MakeForm<T>({ defaultValues, fieldNames }: { defaultValues: T; fieldNames: [] }) {
  const form = useForm<T>({
    defaultValues: defaultValues,
    onSubmit: async (value) => {
      console.log(value);
    }
  });
  return (
    <div>
      {fieldNames.map(name => {
        return (
          <form.Field
            name={name}
            children={(field) => {
              return (
                <div>
                  <label className="font-bold">Group Name:</label>
                  <input
                    placeholder="Gruppenname"
                    id={field.name.toString()}
                    name={field.name.toString()}
                    value={String(field.state.value)}
                    onBlur={field.handleBlur}
                    onChange={e => field.handleChange(e.target.value)}
                  />
                </div>
              );
            }}
          ></form.Field>
        )
      })}

    </div>
  );
}

function FormField() {
  return (<div>
    <form
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
    >
  </div>)
}
