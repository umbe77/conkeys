package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"conkeys/storage"
	"conkeys/utility"
)

func GetUsers(u storage.UserStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Query("filter")
		users, err := u.GetUsers(query)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(users)
	}
}

func GetUser(u storage.UserStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userName := c.Params("username")
		usr, err := u.Get(userName)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(usr)
	}
}

func AddUser(u storage.UserStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var usr storage.User
		if err := c.BodyParser(&usr); err != nil {
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err := u.Add(usr)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"message": "User Created",
		})
	}
}

func UpdateUser(u storage.UserStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var usr storage.User
		if err := c.BodyParser(&usr); err != nil {
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err := u.Update(usr)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{

			"message": "User Updated",
		})
	}
}

func SetPassword(u storage.UserStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userName := c.Params("username")
		if password, ok := c.GetReqHeaders()["X-Pwd"]; ok {

			fmt.Printf("pwd: %s\n", password)
			err := u.SetPassword(userName, utility.EncondePassword(password))
			if err != nil {
				// TODO: Set logging
				return err
			}
			c.Status(fiber.StatusNoContent)
			return nil
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Cannot save an empty password",
		})
	}
}

func DeleteUser(u storage.UserStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userName := c.Params("username")

		err := u.Delete(userName)
		if err != nil {
			// TODO: Set logging
			return err
		}
		return c.JSON(fiber.Map{
			"message": "User Deleted",
		})
	}
}
