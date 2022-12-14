const { User } = require("./models");

const getLogin = async (req, res, next) => {
  let token = req.cookies.login;
  if (
    req.headers.authorization &&
    req.headers.authorization.startsWith("Bearer")
  ) {
    token = req.headers.authorization.split(" ")[1];
  } else if (req.cookies?.auth) {
    token = req.cookies?.auth;
  }
  if (!token) return next({ status: 401, message: "Not Authorized" });
  const user = await User.findByToken(token);
  console.log(user);
  try {
    if (user) {
      if (user.email_verified) {
        req.user = user.toJSON();
        next();
      } else {
        next({ message: "Email not verified", status: 401 });
      }
    } else {
      next({ message: "User not found", status: 400 });
    }
  } catch (err) {
    console.log(err);
    if (err === "Email not verified") {
      res.redirect("/");
    } else if (err.name === "MongooseError") {
      next({ message: "Server Error", status: 500 });
    }
  }
};

const checkMan = (req, res, next) => {
  if (req.user.type === "manager") {
    next();
  } else {
    next("Unauthorized Request");
  }
};

const loginFlag = (req, res, next) => {
  const cookie = req.cookies.login;
  if (cookie) {
    res.redirect("/home");
  } else {
    next();
  }
};

const checkAdmin = (req, res, next) => {
  if (req.user.type === "admin") {
    next();
  } else {
    res.status(403).redirect("/admin/login");
  }
};

const errorHandler = (err, req, res, next) => {
  console.log(err);
  res.status(err.status).json({ message: err.message });
};

module.exports = { getLogin, loginFlag, errorHandler, checkMan, checkAdmin };
