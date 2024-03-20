package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.dockerbuild.Random;
import com.pulumi.dockerbuild.RandomArgs;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        var myRandomResource = new Random("myRandomResource", RandomArgs.builder()        
            .length(24)
            .build());

        ctx.export("value", myRandomResource.result());
    }
}
