import { CredentialResponse, GoogleLogin } from "@react-oauth/google";
import { useUserStore } from "../lib/stores/stores";
import { useNavigate } from "react-router-dom";

export const LoginPage = () => {
  const { user, setUser } = useUserStore((state) => state);
  const navigate = useNavigate();

  if (user.userId > -1) {
    navigate("/dashboard");
  }
  const mocklogin = (driver: boolean) => {
    setUser({
      userId: 1,
      username: "TestUser",
      age: 17,
      email: "test@test.com",
      driver: driver,
      residence: "Graz",
    });
  };

  const onSuccessfulLogin = (res: CredentialResponse) => {
    alert(res);
  };
  const errorOnLogin = () => {
    alert("Login Failed");
  };
  return (
    <div className="flex flex-col justify-center items-center h-screen ">
      <div>
        <h1>RideShare</h1>
        <GoogleLogin
          onSuccess={(res) => onSuccessfulLogin(res)}
          onError={errorOnLogin}
        ></GoogleLogin>
      </div>
      <button
        onClick={() => mocklogin(true)}
        className="border my-1 px-2 border-neutral-400 rounded"
      >
        Mock Login (Driver)
      </button>
      <button
        onClick={() => mocklogin(false)}
        className="border my-1 px-2 border-neutral-400 rounded"
      >
        Mock Login (Not Driver)
      </button>
    </div>
  );
};
