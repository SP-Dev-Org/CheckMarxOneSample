package com.example.api;

import java.sql.*;
import javax.servlet.http.*;
import java.io.*;

public class UserController extends HttpServlet {

    // SQL Injection vulnerability
    public void getUserById(HttpServletRequest request, HttpServletResponse response)
            throws SQLException, IOException {
        String userId = request.getParameter("id");
        Connection conn = DriverManager.getConnection("jdbc:mysql://localhost/mydb", "root", "password");

        // VULNERABLE: Direct concatenation of user input into SQL query
        String query = "SELECT * FROM users WHERE id = " + userId;
        Statement stmt = conn.createStatement();
        ResultSet rs = stmt.executeQuery(query);

        response.getWriter().println("User found: " + rs.getString("name"));
    }

    // Command Injection vulnerability
    public void executeCommand(HttpServletRequest request, HttpServletResponse response)
            throws IOException {
        String fileName = request.getParameter("file");

        // VULNERABLE: User input directly in system command
        Runtime.getRuntime().exec("cat /var/log/" + fileName);

        response.getWriter().println("File content displayed");
    }

    // Path Traversal vulnerability
    public void readFile(HttpServletRequest request, HttpServletResponse response)
            throws IOException {
        String filename = request.getParameter("filename");

        // VULNERABLE: No validation of file path
        File file = new File("/app/data/" + filename);
        BufferedReader reader = new BufferedReader(new FileReader(file));

        String line;
        while ((line = reader.readLine()) != null) {
            response.getWriter().println(line);
        }
        reader.close();
    }

    // XSS vulnerability
    public void displayMessage(HttpServletRequest request, HttpServletResponse response)
            throws IOException {
        String message = request.getParameter("msg");

        // VULNERABLE: Unescaped user input in HTML response
        response.getWriter().println("<html><body>");
        response.getWriter().println("<h1>" + message + "</h1>");
        response.getWriter().println("</body></html>");
    }

    // Hardcoded credentials
    private void connectToDatabase() throws SQLException {
        // VULNERABLE: Hardcoded password
        String dbUrl = "jdbc:mysql://localhost/mydb";
        String username = "admin";
        String password = "SuperSecret123!";

        Connection conn = DriverManager.getConnection(dbUrl, username, password);
    }
}
